package krkstops

import (
	"testing"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	redi "github.com/gomodule/redigo/redis"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/matryer/is"
)

func TestCacheDepartures(t *testing.T) {
	if !testing.Short() {
		is := is.New(t)
		mustClearCache(is)
		depsExpire = time.Millisecond * 10
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
		mustDepsNotBeCached(is, &testStop)
		err := cache.cacheDepartures(&testDepartures, &testStop)
		is.NoErr(err)
		cachedDeps, err := cache.getCachedDepartures(&testStop)
		is.NoErr(err)
		opts := cmpopts.IgnoreUnexported(*testDepartures.Departures[0])
		for i := range cachedDeps.Departures {
			if diff := cmp.Diff(*cachedDeps.Departures[i], *testDepartures.Departures[i], opts); diff != "" {
				t.Errorf(diff)
			}
		}
		time.Sleep(time.Millisecond * 11)
		mustDepsNotBeCached(is, &testStop)
	}
}

func mustDepsNotBeCached(is *is.I, s *pb.Stop) {
	_, err := cache.getCachedDepartures(s)
	is.Equal(err, redi.ErrNil)
}
