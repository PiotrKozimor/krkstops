package cache

import (
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/gomodule/redigo/redis"
)

const (
	AirlyExpire = time.Minute * 10
	airlyPrefix = "airly-"
)

func getAirlyKey(i *pb.Installation) string {
	return airlyPrefix + string(i.Id)
}
func (c *Cache) Airly(airly *pb.Airly, installation *pb.Installation, exp time.Duration, conn redis.Conn) error {
	return c.message(getAirlyKey(installation), airly, exp, conn)
}

func (c *Cache) GetAirly(installation *pb.Installation, conn redis.Conn) (airly *pb.Airly, err error) {
	airly = &pb.Airly{}
	err = c.get(getAirlyKey(installation), airly, conn)
	return
}
