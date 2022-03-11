package krkstops

import (
	"context"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/ttss"
	"github.com/gomodule/redigo/redis"
)

type uniqueStops map[uint32]string

func (s *Score) Update() error {
	_, err := s.conn.Do("DEL", TO_SCORE)
	if err != nil {
		return err
	}
	busStops, err := ttss.BusEndpoint.GetAllStops()
	if err != nil {
		return err
	}
	s.fillIdSet(BUS, busStops)
	tramStops, err := ttss.TramEndpoint.GetAllStops()
	if err != nil {
		return err
	}
	s.fillIdSet(TRAM, tramStops)
	uniqueStops := make(map[uint32]string, len(busStops))
	for i := range busStops {
		uniqueStops[busStops[i].Id] = busStops[i].Name
	}
	for i := range tramStops {
		uniqueStops[tramStops[i].Id] = tramStops[i].Name
	}
	err = s.fillNamesHash(uniqueStops)
	if err != nil {
		return err
	}
	err = s.fillSuggestions(uniqueStops)
	if err != nil {
		return err
	}
	return s.finishUpdate()
}

func (s *Score) fillIdSet(key string, stops []pb.Stop) error {
	ids := make([]interface{}, len(stops))
	for i := range stops {
		ids[i] = stops[i].Id
	}
	args := append(
		[]interface{}{getTmpKey(key)},
		ids...,
	)
	_, err := s.conn.Do("SADD", args...)
	return err
}

func (s *Score) fillNamesHash(stops uniqueStops) error {
	args := make([]interface{}, 2*len(stops))
	i := 0
	for id, name := range stops {
		args[i] = id
		args[i+1] = name
		i += 2
	}
	_, err := s.conn.Do("HSET", append([]interface{}{getTmpKey(NAMES)}, args...)...)
	return err
}

func (s *Score) fillSuggestions(stops uniqueStops) error {
	for id, name := range stops {
		score, err := redis.Float64(s.conn.Do("HGET", SCORES, id))
		if err != nil {
			if err == redis.ErrNil {
				score = 1.0
				_, err := s.conn.Do("SADD", TO_SCORE, id)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		err = addSuggestion(s.sugTmp, &pb.Stop{
			Name: name,
			Id:   id,
		}, score)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Score) finishUpdate() error {
	commands := []string{
		BUS,
		TRAM,
		NAMES,
		SUG,
	}
	for _, cmd := range commands {
		err := s.conn.Send("RENAME", getTmpKey(cmd), cmd)
		if err != nil {
			return err
		}
	}
	err := s.conn.Flush()
	if err != nil {
		return err
	}
	for i := 0; i < len(commands); i++ {
		_, err = s.conn.Receive()
		if err != nil {
			return err
		}
	}
	_, err = s.redis.BgSave(context.Background()).Result()
	return err
}
