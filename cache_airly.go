package krkstops

import (
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
)

var airlyExpire = time.Minute * 10

const airlyPrefix = "airly-"

func getAirlyKey(i *pb.Installation) string {
	return airlyPrefix + string(i.Id)
}
func (c *Cache) cacheAirly(airly *pb.Airly, installation *pb.Installation) error {
	return c.message(getAirlyKey(installation), airly, airlyExpire)
}

func (c *Cache) getCachedAirly(installation *pb.Installation) (airly *pb.Airly, err error) {
	airly = &pb.Airly{}
	err = c.get(getAirlyKey(installation), airly)
	return
}
