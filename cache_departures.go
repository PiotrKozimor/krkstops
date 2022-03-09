package krkstops

import (
	"strconv"
	"time"

	"context"

	"github.com/PiotrKozimor/krkstops/pb"
)

var depsExpire = time.Second * 15

const depsPrefix = "deps-"

func getDeparturesKey(d *pb.Stop) string {
	return depsPrefix + strconv.Itoa(int(d.Id))
}

func (c *Cache) departures(deps *pb.Departures, stop *pb.Stop) (err error) {
	return c.message(getDeparturesKey(stop), deps, depsExpire)
}

func (c *Cache) getDepartures(ctx context.Context, stop *pb.Stop) (departures *pb.Departures, err error) {
	departures = &pb.Departures{}
	err = c.get(getDeparturesKey(stop), departures)
	if err != nil {
		return
	}
	var ttl time.Duration
	ttl, err = c.redis.TTL(ctx, depsPrefix+strconv.Itoa(int(stop.Id))).Result()
	if err != nil {
		return
	}
	livedFor := int32(depsExpire.Seconds() - ttl.Seconds())
	for index := range departures.Departures {
		if departures.Departures[index].RelativeTime != 0 {
			departures.Departures[index].RelativeTime -= livedFor
		}
	}
	return
}
