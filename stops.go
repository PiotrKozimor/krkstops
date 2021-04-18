package krkstops

import (
	"context"
	"strconv"
	"strings"

	"github.com/PiotrKozimor/krkstops/ttss"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type Clients struct {
	Redis              *redis.Client
	RedisAutocompleter *redisearch.Autocompleter
	ctx                context.Context
}

var ctx = context.Background()

// CompareStops and return new and old. New stops are stored in 'stops.tmp' set.
func (c *Clients) CompareStops(stops ttss.Stops) (newStops ttss.Stops, oldStops ttss.Stops, err error) {
	newStops = make(ttss.Stops)
	oldStops = make(ttss.Stops)
	pipe := c.Redis.TxPipeline()
	for ShortName := range stops {
		pipe.SAdd(ctx, "stops.tmp", ShortName)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, nil, err
	}
	newStopIds := c.Redis.SDiff(ctx, "stops.tmp", "stops").Val()
	oldStopIds := c.Redis.SDiff(ctx, "stops", "stops.tmp").Val()
	for _, ShortName := range newStopIds {
		id, err := strconv.Atoi(ShortName)
		if err != nil {
			logrus.Error(err)
		}
		newStops[uint32(id)] = stops[uint32(id)]
	}
	for _, ShortName := range oldStopIds {
		id, err := strconv.Atoi(ShortName)
		if err != nil {
			logrus.Error(err)
		}
		oldStops[uint32(id)] = stops[uint32(id)]
	}
	return newStops, oldStops, nil
}

// UpdateSuggestionsAndRedis in Redisearch engine
func (c *Clients) UpdateSuggestionsAndRedis(newStops ttss.Stops, oldStops ttss.Stops) error {
	pipe := c.Redis.TxPipeline()
	for shortName, name := range newStops {
		pipe.SAdd(ctx, "stops.new", shortName)
		pipe.Set(ctx, strconv.Itoa(int(shortName)), name, -1)
		c.AddSuggestion(strconv.Itoa(int(shortName)), name, 1.0)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	for _, name := range oldStops {
		err := c.RedisAutocompleter.DeleteTerms(redisearch.Suggestion{Term: name})
		if err != nil {
			return err
		}
	}
	_, err = c.Redis.Rename(ctx, "stops.tmp", "stops").Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *Clients) AddSuggestion(shortName, name string, score float64) error {
	splitted := strings.Split(name, " ")
	for index, value := range splitted {
		if value == "(nż)" {
			splitted = append(splitted[:index], splitted[index+1:]...)
		}
	}
	for i := 0; i < len(splitted); i++ {
		swapped := append(splitted[i:], splitted[:i]...)
		term := strings.Join(swapped, " ")
		err := c.RedisAutocompleter.AddTerms(redisearch.Suggestion{Term: term, Score: score, Payload: shortName})
		if err != nil {
			return err
		}
	}
	return nil
}
