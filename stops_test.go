package krkstops

import (
	"context"
	"testing"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/go-redis/redis/v8"
	redi "github.com/gomodule/redigo/redis"
	"github.com/matryer/is"
)

func TestUpdate(t *testing.T) {
	// to avoid err: ERR Background save already in progress
	time.Sleep(time.Millisecond * 100)
	is := is.New(t)
	mustClearCache(is)
	err := score.Update()
	is.NoErr(err)
	busStops, err := score.redis.SMembers(ctx, BUS).Result()
	is.NoErr(err)
	mustHaveSameElements(is, busStops, []string{
		"610",
		"81",
	})
	tramStops, err := score.redis.SMembers(ctx, TRAM).Result()
	is.NoErr(err)
	mustHaveSameElements(is, tramStops, []string{
		"610",
	})
	names, err := score.redis.HGetAll(ctx, NAMES).Result()
	is.NoErr(err)
	is.Equal(names, map[string]string{
		"610": "Rondo Matecznego",
		"81":  "Czarnowiejska",
	})
	toScore, err := score.redis.SMembers(ctx, TO_SCORE).Result()
	is.NoErr(err)
	mustHaveSameElements(is, toScore, []string{
		"610",
		"81",
	})
	exists, err := score.redis.Exists(ctx, SCORES).Result()
	is.NoErr(err)
	is.Equal(exists, int64(0))
	stops, err := score.Search(ctx, "cza")
	is.NoErr(err)
	is.Equal(len(stops), 1)
	is.Equal(stops[0], &pb.Stop{
		Id:   81,
		Name: "Czarnowiejska",
	})
	stops, err = score.Search(ctx, "ma")
	is.NoErr(err)
	is.Equal(len(stops), 1)
	is.Equal(stops[0], &pb.Stop{
		Id:   610,
		Name: "Rondo Matecznego",
	})
}

// mustSeachEqal(is *is.I, )

func BenchmarkRedis(b *testing.B) {
	cli := redis.NewClient(&redis.Options{})
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
	c, err := redi.Dial("tcp", ":6379")
	if err != nil {
		b.Fatal(err)
	}
	var res string
	for i := 0; i < b.N; i++ {
		_, err := c.Do("SET", "FOO", "BAR")
		if err != nil {
			b.Fatal(err)
		}
		res, err = redi.String(c.Do("GET", "FOO"))
		if err != nil {
			b.Fatal(err)
		}
	}
	res += "1"
}
