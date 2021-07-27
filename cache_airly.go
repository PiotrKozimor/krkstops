package krkstops

import (
	"log"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"google.golang.org/protobuf/proto"
)

var AirlyExpire = time.Minute * 10

const airlyPrefix = "airly-"

func getAirlyKey(i *pb.Installation) string {
	return airlyPrefix + string(i.Id)
}
func (c *Cache) cacheAirly(airly *pb.Airly, installation *pb.Installation) (err error) {
	rawAirly, err := proto.Marshal(airly)
	if err != nil {
		return err
	}
	_, err = c.redis.Set(ctx, getAirlyKey(installation), string(rawAirly), AirlyExpire).Result()
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (c *Cache) getCachedAirly(installation *pb.Installation) (*pb.Airly, error) {
	var airly pb.Airly
	airlyRaw, err := c.redis.Get(ctx, getAirlyKey(installation)).Result()
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal([]byte(airlyRaw), &airly)
	return &airly, err
}
