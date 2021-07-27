package krkstops

import (
	"context"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v8"
	redi "github.com/gomodule/redigo/redis"
)

const (
	TO_SCORE = "to.score"
	BUS      = "bus"
	TRAM     = "tram"
	NAMES    = "names"
	SUG      = "sug"
	SCORES   = "scores"
)

var ctx = context.Background()

type Cache struct {
	Redis  *redis.Client
	conn   redi.Conn
	sugTmp *redisearch.Autocompleter
	sug    *redisearch.Autocompleter
}

func NewCache(redisURI, sugKey string) (*Cache, error) {
	conn, err := redi.Dial("tcp", redisURI)
	if err != nil {
		return nil, err
	}
	ac := redisearch.NewAutocompleter(redisURI, sugKey)
	acTmp := redisearch.NewAutocompleter(redisURI, getTmpKey(sugKey))
	r := redis.NewClient(&redis.Options{
		Addr: redisURI,
	})
	db := Cache{
		conn:   conn,
		sugTmp: acTmp,
		sug:    ac,
		Redis:  r,
	}
	return &db, err
}

func getTmpKey(k string) string {
	return "tmp." + k
}
