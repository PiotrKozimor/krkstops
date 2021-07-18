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

func TestCacheDepartures(t *testing.T) {
	if !testing.Short() {
		depsExpire = time.Second * 1
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		client := redis.NewClient(
			&redis.Options{
				Addr: "localhost:6380"})
		testStop := pb.Stop{Name: "Nor", ShortName: "45"}
		testDepartures := []pb.Departure{
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
		}
		_, err := client.Del(context.Background(), getDeparturesKey(&testStop)).Result()
		assert.NoError(t, err)
		isCached, err := isDepartureCached(client, &testStop)
		assert.NoError(t, err)
		assert.False(t, isCached)
		err = cacheDepartures(client, testDepartures, &testStop)
		assert.NoError(t, err)
		isCached, err = isDepartureCached(client, &testStop)
		assert.NoError(t, err)
		assert.True(t, isCached)
		cachedDeps, err := getCachedDepartures(client, &testStop)
		assert.NoError(t, err)
		for index, cachedDep := range cachedDeps {
			if diff := cmp.Diff(cachedDep, testDepartures[index], cmpopts.IgnoreUnexported(cachedDep)); diff != "" {
				t.Errorf(diff)
			}
		}
		time.Sleep(time.Millisecond * 1001)
		isCached, err = isDepartureCached(client, &testStop)
		assert.NoError(t, err)
		assert.False(t, isCached)
	}
}
