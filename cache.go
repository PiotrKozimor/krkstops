package krkstops

import (
	"log"
	"time"

	"github.com/PiotrKozimor/krkstops/ttss"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v8"
	redi "github.com/gomodule/redigo/redis"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	TO_SCORE = "to.score"
	BUS      = "bus"
	TRAM     = "tram"
	NAMES    = "names"
	SUG      = "sug"
	SCORES   = "scores"
)

// var ctx = context.Background()

type Cache struct {
	redis  *redis.Client
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
		redis:  r,
	}
	return &db, err
}

func getTmpKey(k string) string {
	return "tmp." + k
}

func (c *Cache) get(key string, val protoreflect.ProtoMessage) error {
	raw, err := redi.Bytes(c.conn.Do("GET", key))
	if err != nil {
		return err
	}
	return proto.Unmarshal(raw, val)
}

func (c *Cache) message(key string, val protoreflect.ProtoMessage, exp time.Duration) error {
	raw, err := proto.Marshal(val)
	if err != nil {
		return err
	}
	_, err = c.conn.Do("SET", key, raw, "PX", exp.Milliseconds())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c *Cache) getEndpoints(ctx context.Context, id uint32) ([]ttss.Endpoint, error) {
	var endp []ttss.Endpoint
	is, err := c.redis.SIsMember(ctx, BUS, id).Result()
	if err != nil {
		return nil, err
	}
	if is {
		endp = append(endp, ttss.BusEndpoint)
	}
	is, err = c.redis.SIsMember(ctx, TRAM, id).Result()
	if err != nil {
		return nil, err
	}
	if is {
		endp = append(endp, ttss.TramEndpoint)
	}
	return endp, nil
}
