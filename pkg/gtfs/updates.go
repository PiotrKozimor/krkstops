package gtfs

import (
	"fmt"
	"log"
	"time"

	"github.com/PiotrKozimor/krkstops/pkg/gtfs/realtimepb"
	"google.golang.org/protobuf/proto"
)

type stopUpdateKey struct {
	stopId uint32
	tripId uint32
}

type stopUpdates map[stopUpdateKey]uint32

func ParseFeed(b []byte) (*realtimepb.FeedMessage, error) {
	var updates realtimepb.FeedMessage
	err := proto.Unmarshal(b, &updates)
	return &updates, err
}

func (u *Unmarshaler) parseStopUpdates(feed *realtimepb.FeedMessage) (stopUpdates, error) {
	updates := make(stopUpdates, 1000)

	for i := len(feed.Entity) - 1; i >= 0; i-- {
		ent := feed.Entity[i]
		if tu := ent.TripUpdate; tu != nil {
			tripId, err := u.reduceTripId(tu.GetTrip().GetTripId())
			if err != nil {
				return nil, fmt.Errorf("reduce trip id: %w", err)
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
					updateTime := time.Unix(t, 0).In(location)
					secondsInDay := uint32(updateTime.Hour()*60*60 + updateTime.Minute()*60 + updateTime.Second())

					key := stopUpdateKey{
						stopId: stopId,
						tripId: tripId,
					}

					if seconds, ok := updates[key]; ok {
						log.Printf("duplicate update for stop %d and trip %d: %d", stopId, tripId, seconds)
					}
					updates[key] = secondsInDay
				}
			}
		}
	}
	return updates, nil
}
