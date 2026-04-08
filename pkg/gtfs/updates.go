package gtfs

import (
	"fmt"
	"time"

	"github.com/PiotrKozimor/krkstops/pkg/gtfs/realtimepb"
	"google.golang.org/protobuf/proto"
)

type stopUpdateKey struct {
	stopId uint32
	tripId uint32
}

type StopUpdates map[stopUpdateKey]uint32

func ParseFeed(b []byte) (*realtimepb.FeedMessage, error) {
	var updates realtimepb.FeedMessage
	err := proto.Unmarshal(b, &updates)
	return &updates, err
}

func (u *Unmarshaler) ParseStopUpdates(feed *realtimepb.FeedMessage) (StopUpdates, error) {
	updates := make(StopUpdates, 1000)

	for i := len(feed.Entity) - 1; i >= 0; i-- {
		ent := feed.Entity[i]
		if tu := ent.TripUpdate; tu != nil {
			tripId, ok := u.tripIds[tu.GetTrip().GetTripId()]
			if !ok {
				continue
			}
			for _, stu := range tu.StopTimeUpdate {
				stopId, err := u.reduceStopId(stu.GetStopId())
				if err != nil {
					return nil, fmt.Errorf("reduce trip id: %w", err)
				}
				var t int64
				if stu.Departure != nil {
					t = stu.Departure.GetTime()
				} else if stu.Arrival != nil {
					t = stu.Arrival.GetTime()
				}
				if t != 0 {
					updateTime := time.Unix(t, 0).In(Location)
					secondsInDay := uint32(updateTime.Hour()*60*60 + updateTime.Minute()*60 + updateTime.Second())

					key := stopUpdateKey{
						stopId: stopId,
						tripId: tripId,
					}

					if _, ok := updates[key]; !ok {
						updates[key] = secondsInDay
					} else {
						// TODO bump prometheus metric
					}
				}
			}
		}
	}
	return updates, nil
}
