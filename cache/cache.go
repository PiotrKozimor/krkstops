package cache

import (
	"log"

	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
	"github.com/go-redis/redis/v7"
	"google.golang.org/protobuf/proto"
)

func IsCached(c *redis.Client, stop *pb.Stop) (cached bool, err error) {
	var exist int64
	exist, err = c.Exists(stop.ShortName + "-cache").Result()
	if err != nil {
		return
	}
	if err != nil || exist == 0 {
		cached = false
	} else {
		cached = true
	}
	return
}

func CacheDepartures(c *redis.Client, deps *[]pb.Departure, stop *pb.Stop) (err error) {
	pipe := c.Pipeline()
	executePipe := true
	pipe.Del(stop.ShortName + "-cache")
	rawDeps := make([]interface{}, len(*deps))
	for index, dep := range *deps {
		rawDeps[index], err = proto.Marshal(&dep)
		if err != nil {
			log.Println(err)
			return
		}
	}
	pipe.RPush(stop.ShortName+"-cache", rawDeps...)
	if executePipe {
		_, err = pipe.Exec()
		if err != nil {
			log.Println(err)
			return
		}
	}
	return
}

func GetCachedDepartures(c *redis.Client, stop *pb.Stop) (departures []pb.Departure, err error) {
	countDepartures, err := c.LLen(stop.ShortName + "-cache").Result()
	if err != nil {
		return
	}
	rawDepartures, err := c.LRange(stop.ShortName+"-cache", 0, -1).Result()
	departures = make([]pb.Departure, countDepartures)
	if err != nil {
		return
	}
	for index, rawDep := range rawDepartures {
		err = proto.Unmarshal([]byte(rawDep), &departures[index])
		if err != nil {
			return
		}
	}
	return
}
