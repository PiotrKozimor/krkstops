package cache

import (
	"log"
	"strconv"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

var DepsExpire = time.Second * 15

const DepsPrefix = "deps-"

func getDeparturesKey(d *pb.Stop) string {
	return DepsPrefix + strconv.Itoa(int(d.Id))
}

func IsDepartureCached(c *redis.Client, stop *pb.Stop) (cached bool, err error) {
	var exist int64
	exist, err = c.Exists(ctx, getDeparturesKey(stop)).Result()
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

func CacheDepartures(c *redis.Client, deps []pb.Departure, stop *pb.Stop) (err error) {
	pipe := c.Pipeline()
	executePipe := true
	pipe.Del(ctx, getDeparturesKey(stop))
	rawDeps := make([]interface{}, len(deps))
	for index := range deps {
		rawDeps[index], err = proto.Marshal(&(deps[index]))
		if err != nil {
			log.Println(err)
			return
		}
	}
	pipe.RPush(ctx, getDeparturesKey(stop), rawDeps...)
	pipe.Expire(ctx, getDeparturesKey(stop), DepsExpire)
	if executePipe {
		_, err = pipe.Exec(ctx)
		if err != nil {
			log.Println(err)
			return
		}
	}
	return
}

func GetCachedDepartures(c *redis.Client, stop *pb.Stop) (departures []pb.Departure, err error) {
	countDepartures, err := c.LLen(ctx, getDeparturesKey(stop)).Result()
	if err != nil {
		return
	}
	rawDepartures, err := c.LRange(ctx, getDeparturesKey(stop), 0, -1).Result()
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
