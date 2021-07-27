package krkstops

import (
	"strconv"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/ttss"
	"github.com/RediSearch/redisearch-go/redisearch"
	redi "github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

type uniqueStops map[uint32]string

func (db *Cache) Search(phrase string) ([]*pb.Stop, error) {
	stops, err := db.sug.SuggestOpts(
		phrase, redisearch.SuggestOptions{
			Num:          10,
			Fuzzy:        true,
			WithPayloads: true,
			WithScores:   false,
		})
	if err != nil {
		return nil, err
	}
	stopsP := make([]*pb.Stop, len(stops))
	for i, stop := range stops {
		name, err := db.redis.HGet(ctx, NAMES, stop.Payload).Result()
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(stop.Payload)
		if err != nil {
			logrus.Errorf("failed to parse %s to int", stop.Payload)
		} else {
			stopsP[i] = &pb.Stop{Name: name, Id: uint32(id)}
		}
	}
	return stopsP, nil
}

func (db *Cache) Update() error {
	_, err := db.conn.Do("DEL", TO_SCORE)
	if err != nil {
		return err
	}
	busStops, err := ttss.BusEndpoint.GetAllStops()
	if err != nil {
		return err
	}
	db.fillIdSet(BUS, busStops)
	tramStops, err := ttss.TramEndpoint.GetAllStops()
	if err != nil {
		return err
	}
	db.fillIdSet(TRAM, tramStops)
	uniqueStops := make(map[uint32]string, len(busStops))
	for i := range busStops {
		uniqueStops[busStops[i].Id] = busStops[i].Name
	}
	for i := range tramStops {
		uniqueStops[tramStops[i].Id] = tramStops[i].Name
	}
	err = db.fillNamesHash(uniqueStops)
	if err != nil {
		return err
	}
	err = db.fillSuggestions(uniqueStops)
	if err != nil {
		return err
	}
	return db.finishUpdate()
}

func (db *Cache) fillIdSet(key string, stops []pb.Stop) error {
	ids := make([]interface{}, len(stops))
	for i := range stops {
		ids[i] = stops[i].Id
	}
	args := append(
		[]interface{}{getTmpKey(key)},
		ids...,
	)
	_, err := db.conn.Do("SADD", args...)
	return err
}

func (db *Cache) fillNamesHash(stops uniqueStops) error {
	args := make([]interface{}, 2*len(stops))
	i := 0
	for id, name := range stops {
		args[i] = id
		args[i+1] = name
		i += 2
	}
	_, err := db.conn.Do("HSET", append([]interface{}{getTmpKey(NAMES)}, args...)...)
	return err
}

func (c *Cache) fillSuggestions(stops uniqueStops) error {
	for id, name := range stops {
		score, err := redi.Float64(c.conn.Do("HGET", SCORES, id))
		if err != nil {
			if err == redi.ErrNil {
				score = 1.0
				_, err := c.conn.Do("SADD", TO_SCORE, id)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		err = c.addSuggestion(&pb.Stop{
			Name: name,
			Id:   id,
		}, score)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Cache) finishUpdate() error {
	commands := []string{
		BUS,
		TRAM,
		NAMES,
		SUG,
	}
	for _, cmd := range commands {
		err := db.conn.Send("RENAME", getTmpKey(cmd), cmd)
		if err != nil {
			return err
		}
	}
	err := db.conn.Flush()
	if err != nil {
		return err
	}
	for i := 0; i < len(commands); i++ {
		_, err = db.conn.Receive()
		if err != nil {
			return err
		}
	}
	return nil
}
