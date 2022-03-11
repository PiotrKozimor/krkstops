package cache

import (
	"testing"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/gomodule/redigo/redis"
	"github.com/matryer/is"
)

func mustClearCache(is *is.I, c *Cache) {
	conn := c.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHALL")
	is.NoErr(err)
}

func TestCacheAirly(t *testing.T) {
	is := is.New(t)
	c := NewCache("localhost:6379", SUG)
	mustClearCache(is, c)
	airlyExpireTest := time.Millisecond * 10
	testInst := pb.Installation{Id: 9914}
	testAirly := pb.Airly{
		Caqi:        32,
		Humidity:    45,
		Color:       435,
		Temperature: 45.6,
	}
	conn := c.Pool.Get()
	defer conn.Close()
	mustAirlyNotBeCached(is, &testInst, c, conn)
	err := c.Airly(&testAirly, &testInst, airlyExpireTest, conn)
	is.NoErr(err)
	cachedAirly, err := c.GetAirly(&testInst, conn)
	is.NoErr(err)

	mustAirlyEqual(is, cachedAirly, &testAirly)
	time.Sleep(time.Millisecond * 11)
	mustAirlyNotBeCached(is, &testInst, c, conn)
}

func mustAirlyNotBeCached(is *is.I, i *pb.Installation, c *Cache, conn redis.Conn) {
	_, err := c.GetAirly(i, conn)
	is.Equal(err, redis.ErrNil)
}

func mustAirlyEqual(is *is.I, a1 *pb.Airly, a2 *pb.Airly) {
	is.Equal(a1.Caqi, a2.Caqi)
	is.Equal(a1.Color, a2.Color)
	is.Equal(a1.ColorStr, a2.ColorStr)
	is.Equal(a1.Humidity, a2.Humidity)
	is.Equal(a1.Temperature, a2.Temperature)
}
