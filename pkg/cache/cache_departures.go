package cache

import (
	"strconv"
	"time"

	"context"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/gomodule/redigo/redis"
)

const (
	DepsExpire = time.Second * 15
	depsPrefix = "deps-"
)

func getDeparturesKey(d *pb.Stop) string {
	return depsPrefix + strconv.Itoa(int(d.Id))
}

func (c *Cache) Departures(deps *pb.Departures, stop *pb.Stop, exp time.Duration, conn redis.Conn) (err error) {
	return c.message(getDeparturesKey(stop), deps, exp, conn)
}

func (c *Cache) GetDepartures(ctx context.Context, stop *pb.Stop, exp time.Duration, conn redis.Conn) (departures *pb.Departures, err error) {
	departures = &pb.Departures{}
	err = c.get(getDeparturesKey(stop), departures, conn)
	if err != nil {
		return
	}
	var ttl int
	ttl, err = redis.Int(conn.Do("TTL", depsPrefix+strconv.Itoa(int(stop.Id))))
	if err != nil {
		return
	}
	livedFor := int32(int(exp.Seconds()) - ttl)
	for index := range departures.Departures {
		if departures.Departures[index].RelativeTime != 0 {
			departures.Departures[index].RelativeTime -= livedFor
		}
	}
	return
}
