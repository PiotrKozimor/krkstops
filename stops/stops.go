package stops

import (
	"strings"

	"github.com/PiotrKozimor/krkstops/ttss"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v7"
)

type Clients struct {
	Redis              *redis.Client
	RedisAutocompleter *redisearch.Autocompleter
}

// CompareStops and return new and old. New stops are stored in 'stops.tmp' set.
func (c *Clients) CompareStops(stops ttss.Stops) (newStops ttss.Stops, oldStops ttss.Stops, err error) {
	newStops = make(ttss.Stops)
	oldStops = make(ttss.Stops)
	pipe := c.Redis.TxPipeline()
	for ShortName := range stops {
		pipe.SAdd("stops.tmp", ShortName)
	}
	_, err = pipe.Exec()
	if err != nil {
		return nil, nil, err
	}
	newStopsShortNames := c.Redis.SDiff("stops.tmp", "stops").Val()
	oldStopsShortNames := c.Redis.SDiff("stops", "stops.tmp").Val()
	for _, ShortName := range newStopsShortNames {
		newStops[ShortName] = stops[ShortName]
	}
	for _, ShortName := range oldStopsShortNames {
		oldStops[ShortName] = stops[ShortName]
	}
	return newStops, oldStops, nil
}

// UpdateSuggestionsAndRedis in Redisearch engine
func (c *Clients) UpdateSuggestionsAndRedis(newStops ttss.Stops, oldStops ttss.Stops) error {
	pipe := c.Redis.TxPipeline()
	for shortName, name := range newStops {
		pipe.SAdd("stops.new", shortName)
		pipe.Set(shortName, name, -1)
		c.AddSuggestion(shortName, name, 1.0)
	}
	_, err := pipe.Exec()
	if err != nil {
		return err
	}
	for _, name := range oldStops {
		err := c.RedisAutocompleter.DeleteTerms(redisearch.Suggestion{Term: name})
		if err != nil {
			return err
		}
	}
	_, err = c.Redis.Rename("stops.tmp", "stops").Result()
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
