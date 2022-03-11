package cache

import (
	"strconv"
	"time"

	"context"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
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

type Cache struct {
	// Conn redis.Conn
	Pool redis.Pool
	Sug  *redisearch.Autocompleter
}

func NewCache(redisURI, sugKey string) *Cache {
	pool := redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisURI)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	ac := redisearch.NewAutocompleterFromPool(&pool, sugKey)
	c := Cache{
		Pool: pool,
		Sug:  ac,
	}
	return &c
}

func (c *Cache) Search(ctx context.Context, phrase string) ([]*pb.Stop, error) {
	stops, err := c.Sug.SuggestOpts(
		phrase, redisearch.SuggestOptions{
			Num:          10,
			Fuzzy:        true,
			WithPayloads: true,
			WithScores:   false,
		})
	if err != nil {
		return nil, err
	}
	stopsP := make([]*pb.Stop, len(stops))
	conn := c.Pool.Get()
	defer conn.Close()
	for i, stop := range stops {
		name, err := redis.String(conn.Do("HGET", NAMES, stop.Payload))
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(stop.Payload)
		if err != nil {
			logrus.Errorf("failed to parse %s to int", stop.Payload)
		} else {
			stopsP[i] = &pb.Stop{Name: name, Id: uint32(id)}
		}
	}
	return stopsP, nil
}

func GetTmpKey(k string) string {
	return "tmp." + k
}

func (c *Cache) get(key string, val protoreflect.ProtoMessage, conn redis.Conn) error {
	raw, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return err
	}
	return proto.Unmarshal(raw, val)
}

func (c *Cache) message(key string, val protoreflect.ProtoMessage, exp time.Duration, conn redis.Conn) error {
	raw, err := proto.Marshal(val)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", key, raw, "PX", exp.Milliseconds())
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetEndpoints(ctx context.Context, id uint32) ([]pb.Endpoint, error) {
	var endp []pb.Endpoint
	conn := c.Pool.Get()
	defer conn.Close()
	is, err := redis.Bool(conn.Do("SISMEMBER", BUS, id))
	if err != nil {
		return nil, err
	}
	if is {
		endp = append(endp, pb.Endpoint_BUS)
	}
	is, err = redis.Bool(conn.Do("SISMEMBER", TRAM, id))
	if err != nil {
		return nil, err
	}
	if is {
		endp = append(endp, pb.Endpoint_TRAM)
	}
	return endp, nil
}
