package krkstops

import (
	"log"
	"os"
	"testing"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/matryer/is"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var (
	stopsDB *Cache
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	var err error
	stopsDB, err = NewCache("localhost:6379", SUG)
	handle(err)
}

func mustClearCache(is *is.I) {
	_, err := stopsDB.conn.Do("FLUSHALL")
	is.NoErr(err)
}

func Test_scoreByTotalDepartures(t *testing.T) {
	assert.Equal(t, float64(1), scoreByTotalDepartures(0))
	assert.Equal(t, float64(5), scoreByTotalDepartures(64))
}

func TestScore(t *testing.T) {
	is := is.New(t)
	mustClearCache(is)
	err := stopsDB.Update()
	is.NoErr(err)
	err = stopsDB.Score(make(<-chan os.Signal), mustLocalKrkStopsClient(is))
	is.NoErr(err)
	scores, err := stopsDB.Redis.HGetAll(ctx, SCORES).Result()
	is.NoErr(err)
	is.Equal(scores, map[string]string{
		"610": "1.0",
		"81":  "1.0",
	})
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
