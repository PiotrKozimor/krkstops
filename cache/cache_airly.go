package cache

import (
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
)

const (
	AirlyExpire = time.Minute * 10
	airlyPrefix = "airly-"
)

func getAirlyKey(i *pb.Installation) string {
	return airlyPrefix + string(i.Id)
}
func (c *Cache) Airly(airly *pb.Airly, installation *pb.Installation, exp time.Duration) error {
	return c.message(getAirlyKey(installation), airly, exp)
}

func (c *Cache) GetAirly(installation *pb.Installation) (airly *pb.Airly, err error) {
	airly = &pb.Airly{}
	err = c.get(getAirlyKey(installation), airly)
	return
}
