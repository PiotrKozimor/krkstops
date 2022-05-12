package score

import (
	"context"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/pkg/cache"
	"github.com/PiotrKozimor/krkstops/pkg/ttss"
	"github.com/gomodule/redigo/redis"
)

type uniqueStops map[uint32]string

func (s *Score) Update(bus ttss.Endpointer, tram ttss.Endpointer) error {
	_, err := s.Conn.Do("DEL", cache.TO_SCORE)
	if err != nil {
		return err
	}
	busStops, err := bus.GetAllStops()
	if err != nil {
		return err
	}
	err = s.fillIdSet(cache.BUS, busStops)
	if err != nil {
		return err
	}
	tramStops, err := tram.GetAllStops()
	if err != nil {
		return err
	}
	err = s.fillIdSet(cache.TRAM, tramStops)
	if err != nil {
		return err
	}
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
	args := make([]interface{}, len(stops)+1)
	args[0] = cache.GetTmpKey(key)
	for i := range stops {
		args[i+1] = stops[i].Id
	}
	_, err := s.Conn.Do("SADD", args...)
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
	_, err := s.Conn.Do("HSET", append([]interface{}{cache.GetTmpKey(cache.NAMES)}, args...)...)
	return err
}

func (s *Score) fillSuggestions(stops uniqueStops) error {
	for id, name := range stops {
		score, err := redis.Float64(s.Conn.Do("HGET", cache.SCORES, id))
		if err != nil {
			if err == redis.ErrNil {
				score = 1.0
				_, err := s.Conn.Do("SADD", cache.TO_SCORE, id)
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
		cache.BUS,
		cache.TRAM,
		cache.NAMES,
		cache.SUG,
	}
	for _, cmd := range commands {
		err := s.Conn.Send("RENAME", cache.GetTmpKey(cmd), cmd)
		if err != nil {
			return err
		}
	}
	err := s.Conn.Flush()
	if err != nil {
		return err
	}
	for i := 0; i < len(commands); i++ {
		_, err = s.Conn.Receive()
		if err != nil {
			return err
		}
	}
	_, err = s.redis.BgSave(context.Background()).Result()
	return err
}
