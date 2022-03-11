package cache

import (
	"log"
	"strconv"
	"time"

	"context"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/ttss"
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
	Conn redis.Conn
	Sug  *redisearch.Autocompleter
}

func NewCache(redisURI, sugKey string) (*Cache, error) {
	conn, err := redis.Dial("tcp", redisURI)
	if err != nil {
		return nil, err
	}
	ac := redisearch.NewAutocompleter(redisURI, sugKey)

	c := Cache{
		Conn: conn,
		Sug:  ac,
	}
	return &c, err
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
	for i, stop := range stops {
		name, err := redis.String(c.Conn.Do("HGET", NAMES, stop.Payload))
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

func (c *Cache) get(key string, val protoreflect.ProtoMessage) error {
	raw, err := redis.Bytes(c.Conn.Do("GET", key))
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
	_, err = c.Conn.Do("SET", key, raw, "PX", exp.Milliseconds())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c *Cache) GetEndpoints(ctx context.Context, id uint32) ([]ttss.Endpoint, error) {
	var endp []ttss.Endpoint
	is, err := redis.Bool(c.Conn.Do("SISMEMBER", BUS, id))
	if err != nil {
		return nil, err
	}
	if is {
		endp = append(endp, ttss.BusEndpoint)
	}
	is, err = redis.Bool(c.Conn.Do("SISMEMBER", TRAM, id))
	if err != nil {
		return nil, err
	}
	if is {
		endp = append(endp, ttss.TramEndpoint)
	}
	return endp, nil
}
