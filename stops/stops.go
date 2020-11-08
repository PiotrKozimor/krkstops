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
func (app *Clients) UpdateSuggestionsAndRedis(newStops ttss.Stops, oldStops ttss.Stops) error {
	pipe := app.Redis.TxPipeline()
	for ShortName, name := range newStops {
		pipe.SAdd("stops.new", ShortName)
		pipe.Set(ShortName, name, -1)
		splitted := strings.Split(name, " ")
		for index, value := range splitted {
			if value == "(nż)" {
				splitted = append(splitted[:index], splitted[index+1:]...)
			}
		}
		for i := 0; i < len(splitted); i++ {
			swapped := append(splitted[i:], splitted[:i]...)
			term := strings.Join(swapped, " ")
			err := app.RedisAutocompleter.AddTerms(redisearch.Suggestion{Term: term, Score: 1, Payload: ShortName})
			if err != nil {
				return err
			}
		}
	}
	_, err := pipe.Exec()
	if err != nil {
		return err
	}
	for _, name := range oldStops {
		err := app.RedisAutocompleter.DeleteTerms(redisearch.Suggestion{Term: name})
		if err != nil {
			return err
		}
	}
	_, err = app.Redis.Rename("stops.tmp", "stops").Result()
	if err != nil {
		return err
	}
	return nil
}
