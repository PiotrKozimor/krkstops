package cache

import (
	"log"
	"time"

	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
	"github.com/go-redis/redis/v7"
	"google.golang.org/protobuf/proto"
)

var AirlyExpire = time.Minute * 10

const AirlyPrefix = "airly-"

func IsAirlyCached(c *redis.Client, installation *pb.Installation) (cached bool, err error) {
	var exist int64
	exist, err = c.Exists(AirlyPrefix + string(installation.Id)).Result()
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

func getAirlyKey(i *pb.Installation) string {
	return AirlyPrefix + string(i.Id)
}
func CacheAirly(c *redis.Client, airly *pb.Airly, installation *pb.Installation) (err error) {
	pipe := c.Pipeline()
	executePipe := true
	pipe.Del(getAirlyKey(installation))
	rawAirly, err := proto.Marshal(airly)
	if err != nil {
		return err
	}
	pipe.Set(getAirlyKey(installation), string(rawAirly), 0)
	pipe.Expire(getAirlyKey(installation), AirlyExpire)
	if executePipe {
		_, err = pipe.Exec()
		if err != nil {
			log.Println(err)
			return
		}
	}
	return
}

func GetCachedAirly(c *redis.Client, installation *pb.Installation) (airly pb.Airly, err error) {
	airlyRaw, err := c.Get(getAirlyKey(installation)).Result()
	if err != nil {
		return airly, err
	}
	err = proto.Unmarshal([]byte(airlyRaw), &airly)
	if err != nil {
		return airly, err
	}
	return
}
