package cache

import (
	"log"
	"testing"

	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
	"github.com/go-redis/redis/v7"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCacheDepartures(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379"})
	testStop := pb.Stop{Name: "Nor", ShortName: "45"}
	testDepartures := []pb.Departure{
		pb.Departure{
			Direction:   "Czerwone Maki Żółć",
			PatternText: "52",
			PlannedTime: "21:43",
		},
		pb.Departure{
			Direction:   "Rząka",
			PatternText: "139",
			PlannedTime: "4:32",
		},
	}
	_, err := client.Del(testStop.ShortName + "-cache").Result()
	if err != nil {
		t.Fatal(err)
	}
	isCached, err := IsCached(client, &testStop)
	if err != nil {
		t.Fatal(err)
	} else if isCached != false {
		t.Fatal("Stops cached")
	}
	err = CacheDepartures(client, &testDepartures, &testStop)
	if err != nil {
		t.Fatal(err)
	}
	isCached, err = IsCached(client, &testStop)
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
}
