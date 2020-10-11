package krkstops

import (
	"testing"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v7"
)

func TestUpdateSuggestionsAndRedis(t *testing.T) {
	app := App{
		RedisClient:        redis.NewClient(&redis.Options{}),
		RedisAutocompleter: redisearch.NewAutocompleter("localhost:6379", "bar"),
	}
	newStops := StopsMap{
		"3": "Budzyń Plaża Główna (nż)",
		"4": "Poczta Główna",
		"5": "Plac Wolnica",
		"6": "Plac Bohaterów Getta",
	}
	oldStops := StopsMap{
		"0": "Foo",
	}
	app.UpdateSuggestionsAndRedis(newStops, oldStops)
}
