package krkstops

import (
	"context"
	"testing"

	"github.com/PiotrKozimor/krkstops/pb"
	redi "github.com/go-redis/redis/v8"
	"github.com/gomodule/redigo/redis"
	"github.com/matryer/is"
)

func TestUpdate(t *testing.T) {
	is := is.New(t)
	mustClearCache(is)
	err := cache.Update()
	is.NoErr(err)
	busStops, err := cache.redis.SMembers(ctx, BUS).Result()
	is.NoErr(err)
	mustHaveSameElements(is, busStops, []string{
		"610",
		"81",
	})
	tramStops, err := cache.redis.SMembers(ctx, TRAM).Result()
	is.NoErr(err)
	mustHaveSameElements(is, tramStops, []string{
		"610",
	})
	names, err := cache.redis.HGetAll(ctx, NAMES).Result()
	is.NoErr(err)
	is.Equal(names, map[string]string{
		"610": "Rondo Matecznego",
		"81":  "Czarnowiejska",
	})
	toScore, err := cache.redis.SMembers(ctx, TO_SCORE).Result()
	is.NoErr(err)
	mustHaveSameElements(is, toScore, []string{
		"610",
		"81",
	})
	exists, err := cache.redis.Exists(ctx, SCORES).Result()
	is.NoErr(err)
	is.Equal(exists, int64(0))
	stops, err := cache.Search("cza")
	is.NoErr(err)
	is.Equal(len(stops), 1)
	is.Equal(stops[0], &pb.Stop{
		Id:   81,
		Name: "Czarnowiejska",
	})
	stops, err = cache.Search("ma")
	is.NoErr(err)
	is.Equal(len(stops), 1)
	is.Equal(stops[0], &pb.Stop{
		Id:   610,
		Name: "Rondo Matecznego",
	})
}

// mustSeachEqal(is *is.I, )

func BenchmarkRedis(b *testing.B) {
	cli := redi.NewClient(&redi.Options{})
	ctx := context.Background()
	var res string
	for i := 0; i < b.N; i++ {
		_, err := cli.Set(ctx, "FOO", "BAR", 0).Result()
		if err != nil {
			b.Fatal(err)
		}
		res, err = cli.Get(ctx, "FOO").Result()
		if err != nil {
			b.Fatal(err)
		}
	}
	res += "1"
}

func BenchmarkRedigo(b *testing.B) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		b.Fatal(err)
	}
	var res string
	for i := 0; i < b.N; i++ {
		_, err := c.Do("SET", "FOO", "BAR")
		if err != nil {
			b.Fatal(err)
		}
		res, err = redis.String(c.Do("GET", "FOO"))
		if err != nil {
			b.Fatal(err)
		}
	}
	res += "1"
}
