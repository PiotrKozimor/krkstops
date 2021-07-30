package krkstops

import (
	"log"
	"os"
	"testing"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/matryer/is"
	"google.golang.org/grpc"
)

var (
	cache *Cache
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	var err error
	cache, err = NewCache("localhost:6379", SUG)
	handle(err)
}

func mustClearCache(is *is.I) {
	_, err := cache.conn.Do("FLUSHALL")
	is.NoErr(err)
}

func Test_scoreByTotalDepartures(t *testing.T) {
	is := is.New(t)
	is.Equal(float64(1), scoreByTotalDepartures(0))
	is.Equal(float64(5), scoreByTotalDepartures(64))
}

func TestScore(t *testing.T) {
	is := is.New(t)
	mustClearCache(is)
	err := cache.Update()
	is.NoErr(err)
	err = cache.Score(ctx, make(<-chan os.Signal), mustLocalKrkStopsClient(is), 0)
	is.NoErr(err)
	scores, err := cache.redis.HGetAll(ctx, SCORES).Result()
	is.NoErr(err)
	is.Equal(scores, map[string]string{
		"610": "3.29128784747792",
		"81":  "2.6583123951777",
	})
	stops, err := cache.sug.SuggestOpts("mat", redisearch.SuggestOptions{
		Num:          10,
		Fuzzy:        true,
		WithPayloads: true,
		WithScores:   true,
	})
	is.NoErr(err)
	is.Equal(len(stops), 1)
	is.Equal(stops[0].Payload, "610")
	is.Equal(stops[0].Score, 0.87963366508483887)
	is.Equal(stops[0].Term, "Matecznego Rondo")
}

func mustLocalKrkStopsClient(is *is.I) pb.KrkStopsClient {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	is.NoErr(err)
	return pb.NewKrkStopsClient(conn)
}

func mustHaveSameElements(is *is.I, a1, a2 []string) {
	is.Equal(len(a1), len(a2))
	m1 := make(map[string]bool, len(a1))
	m2 := make(map[string]bool, len(a2))
	for _, s := range a1 {
		m1[s] = true
	}
	for _, s := range a2 {
		m2[s] = true
	}
	is.Equal(m1, m2)
}
