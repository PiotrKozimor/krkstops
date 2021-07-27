package krkstops

import (
	"testing"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/go-redis/redis/v8"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/matryer/is"
)

func TestCacheAirly(t *testing.T) {
	is := is.New(t)
	mustClearCache(is)
	AirlyExpire = time.Millisecond * 10
	testInst := pb.Installation{Id: 9914}
	testAirly := pb.Airly{
		Caqi:        32,
		Humidity:    45,
		Color:       435,
		Temperature: 45.6,
	}
	mustNotBeCached(is, &testInst)
	err := cache.cacheAirly(&testAirly, &testInst)
	is.NoErr(err)
	cachedAirly, err := cache.getCachedAirly(&testInst)
	is.NoErr(err)
	if diff := cmp.Diff(*cachedAirly, testAirly, cmpopts.IgnoreUnexported(*cachedAirly)); diff != "" {
		t.Errorf(diff)
	}
	time.Sleep(time.Millisecond * 11)
	mustNotBeCached(is, &testInst)
}

func mustNotBeCached(is *is.I, i *pb.Installation) {
	_, err := cache.getCachedAirly(i)
	is.Equal(err, redis.Nil)
}
