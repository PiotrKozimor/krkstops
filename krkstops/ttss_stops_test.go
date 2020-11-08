package krkstops

import (
	"net/http"
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

func TestGetAllStops(t *testing.T) {
	app := App{
		HTTPClient: http.DefaultClient,
	}
	maps, err := app.GetAllStops()
	if err != nil {
		t.Fatal(err)
	}
	if len(maps) < 300 {
		t.Error("Got less than 300 stops!")
	}
}
