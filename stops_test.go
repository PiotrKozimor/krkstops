package krkstops

import (
	"testing"

	"github.com/PiotrKozimor/krkstops/ttss"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v8"
)

func TestUpdateSuggestionsAndRedis(t *testing.T) {
	app := Clients{
		Redis:              redis.NewClient(&redis.Options{}),
		RedisAutocompleter: redisearch.NewAutocompleter("localhost:6379", "bar"),
	}
	newStops := ttss.Stops{
		3: "Budzyń Plaża Główna (nż)",
		4: "Poczta Główna",
		5: "Plac Wolnica",
		6: "Plac Bohaterów Getta",
	}
	oldStops := ttss.Stops{
		0: "Foo",
	}
	app.UpdateSuggestionsAndRedis(newStops, oldStops)
}
