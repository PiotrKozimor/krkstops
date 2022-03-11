package cache

import (
	"testing"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/gomodule/redigo/redis"
	"github.com/matryer/is"
)

func mustClearCache(is *is.I, c *Cache) {
	_, err := c.Conn.Do("FLUSHALL")
	is.NoErr(err)
}

func TestCacheAirly(t *testing.T) {
	is := is.New(t)
	c, err := NewCache("localhost:6379", SUG)
	is.NoErr(err)
	mustClearCache(is, c)
	airlyExpireTest := time.Millisecond * 10
	testInst := pb.Installation{Id: 9914}
	testAirly := pb.Airly{
		Caqi:        32,
		Humidity:    45,
		Color:       435,
		Temperature: 45.6,
	}
	mustAirlyNotBeCached(is, &testInst, c)
	err = c.Airly(&testAirly, &testInst, airlyExpireTest)
	is.NoErr(err)
	cachedAirly, err := c.GetAirly(&testInst)
	is.NoErr(err)

	mustAirlyEqual(is, cachedAirly, &testAirly)
	time.Sleep(time.Millisecond * 11)
	mustAirlyNotBeCached(is, &testInst, c)
}

func mustAirlyNotBeCached(is *is.I, i *pb.Installation, c *Cache) {
	_, err := c.GetAirly(i)
	is.Equal(err, redis.ErrNil)
}

func mustAirlyEqual(is *is.I, a1 *pb.Airly, a2 *pb.Airly) {
	is.Equal(a1.Caqi, a2.Caqi)
	is.Equal(a1.Color, a2.Color)
	is.Equal(a1.ColorStr, a2.ColorStr)
	is.Equal(a1.Humidity, a2.Humidity)
	is.Equal(a1.Temperature, a2.Temperature)
}
