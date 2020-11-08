package cache

import (
	"log"
	"testing"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/go-redis/redis/v7"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCacheDepartures(t *testing.T) {
	DepsExpire = time.Second * 1
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379"})
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
	_, err := client.Del(getDeparturesKey(&testStop)).Result()
	if err != nil {
		t.Fatal(err)
	}
	isCached, err := IsDepartureCached(client, &testStop)
	if err != nil {
		t.Fatal(err)
	} else if isCached != false {
		t.Fatal("Stops cached")
	}
	err = CacheDepartures(client, &testDepartures, &testStop)
	if err != nil {
		t.Fatal(err)
	}
	isCached, err = IsDepartureCached(client, &testStop)
	if err != nil {
		t.Fatal(err)
	} else if isCached != true {
		t.Fatal("Stops not cached")
	}
	cachedDeps, err := GetCachedDepartures(client, &testStop)
	for index, cachedDep := range cachedDeps {
		if diff := cmp.Diff(cachedDep, testDepartures[index], cmpopts.IgnoreUnexported(cachedDep)); diff != "" {
			t.Errorf(diff)
		}
	}
	time.Sleep(time.Millisecond * 1001)
	isCached, err = IsDepartureCached(client, &testStop)
	if err != nil {
		t.Fatal(err)
	} else if isCached != false {
		t.Fatal("Stops not expired")
	}
}
