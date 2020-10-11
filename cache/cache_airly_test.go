package cache

import (
	"log"
	"testing"
	"time"

	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
	"github.com/go-redis/redis/v7"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCacheAirly(t *testing.T) {
	AirlyExpire = time.Second * 1
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379"})
	testInst := pb.Installation{Id: 9914}
	testAirly := pb.Airly{
		Caqi:        32,
		Humidity:    45,
		Color:       "435",
		Temperature: 45.6,
	}
	_, err := client.Del(getAirlyKey(&testInst)).Result()
	if err != nil {
		t.Fatal(err)
	}
	isCached, err := IsAirlyCached(client, &testInst)
	if err != nil {
		t.Fatal(err)
	} else if isCached != false {
		t.Fatal("Instalation cached")
	}
	err = CacheAirly(client, &testAirly, &testInst)
	if err != nil {
		t.Fatal(err)
	}
	isCached, err = IsAirlyCached(client, &testInst)
	if err != nil {
		t.Fatal(err)
	} else if isCached != true {
		t.Fatal("Installation not cached")
	}
	cachedAirly, err := GetCachedAirly(client, &testInst)
	if diff := cmp.Diff(cachedAirly, testAirly, cmpopts.IgnoreUnexported(cachedAirly)); diff != "" {
		t.Errorf(diff)
	}

	time.Sleep(time.Millisecond * 1001)
	isCached, err = IsAirlyCached(client, &testInst)
	if err != nil {
		t.Fatal(err)
	} else if isCached != false {
		t.Fatal("Installation not expired")
	}
}
