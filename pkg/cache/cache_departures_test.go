package cache

import (
	"context"
	"testing"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/gomodule/redigo/redis"
	"github.com/matryer/is"
)

func TestCacheDepartures(t *testing.T) {
	if !testing.Short() {
		is := is.New(t)
		c := NewCache("localhost:6379", SUG)
		mustClearCache(is, c)
		depsExpireTest := time.Millisecond * 10
		testStop := pb.Stop{Name: "Nor", ShortName: "45"}
		testDepartures := pb.Departures{
			Departures: []*pb.Departure{
				{
					Direction:   "Czerwone Maki Żółć",
					PatternText: "52",
					PlannedTime: "21:43",
				},
				{
					Direction:   "Rząka",
					PatternText: "139",
					PlannedTime: "4:32",
				},
			},
		}
		conn := c.Pool.Get()
		defer conn.Close()
		mustDepsNotBeCached(is, &testStop, c, conn)
		err := c.Departures(&testDepartures, &testStop, depsExpireTest, conn)
		is.NoErr(err)
		cachedDeps, err := c.GetDepartures(context.Background(), &testStop, depsExpireTest, conn)
		is.NoErr(err)
		for i := range cachedDeps.Departures {
			mustDepsEqual(is, cachedDeps.Departures[i], testDepartures.Departures[i])
		}
		time.Sleep(time.Millisecond * 11)
		mustDepsNotBeCached(is, &testStop, c, conn)
	}
}

func mustDepsNotBeCached(is *is.I, s *pb.Stop, c *Cache, conn redis.Conn) {
	depsExpireTest := time.Millisecond * 10
	_, err := c.GetDepartures(context.Background(), s, depsExpireTest, conn)
	is.Equal(err, redis.ErrNil)
}

func mustDepsEqual(is *is.I, d1 *pb.Departure, d2 *pb.Departure) {
	is.Equal(d1.Color, d2.Color)
	is.Equal(d1.Direction, d2.Direction)
	is.Equal(d1.PatternText, d2.PatternText)
	is.Equal(d1.PlannedTime, d2.PlannedTime)
	is.Equal(d1.Predicted, d2.Predicted)
	is.Equal(d1.RelativeTime, d2.RelativeTime)
	is.Equal(d1.RelativeTimeParsed, d2.RelativeTimeParsed)
	is.Equal(d1.Type, d2.Type)
}
