package krkstops

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/go-redis/redis/v8"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestCacheAirly(t *testing.T) {
	if !testing.Short() {
		AirlyExpire = time.Second * 1
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		client := redis.NewClient(
			&redis.Options{
				Addr: "localhost:6380"})
		testInst := pb.Installation{Id: 9914}
		testAirly := pb.Airly{
			Caqi:        32,
			Humidity:    45,
			Color:       435,
			Temperature: 45.6,
		}
		_, err := client.Del(context.Background(), getAirlyKey(&testInst)).Result()
		assert.NoError(t, err)
		isCached, err := isAirlyCached(client, &testInst)
		assert.NoError(t, err)
		assert.False(t, isCached)
		err = cacheAirly(client, &testAirly, &testInst)
		assert.NoError(t, err)
		isCached, err = isAirlyCached(client, &testInst)
		assert.NoError(t, err)
		assert.True(t, isCached)
		cachedAirly, err := getCachedAirly(client, &testInst)
		assert.NoError(t, err)
		if diff := cmp.Diff(*cachedAirly, testAirly, cmpopts.IgnoreUnexported(*cachedAirly)); diff != "" {
			t.Errorf(diff)
		}
		time.Sleep(time.Millisecond * 1001)
		isCached, err = isAirlyCached(client, &testInst)
		assert.NoError(t, err)
		assert.False(t, isCached)
	}
}
