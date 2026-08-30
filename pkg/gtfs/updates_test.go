package gtfs

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.cozymore.dev/krkstops/pkg/gtfs/realtimepb"
)

func TestFeed(t *testing.T) {

	for _, source := range []string{
		"A",
		"M",
		"T",
	} {
		b, err := os.ReadFile(fmt.Sprintf("testdata/TripUpdates_%s.pb", source))
		require.NoError(t, err)
		feed, err := ParseFeed(b)
		require.NoError(t, err)

		assertCommon := func(ent *realtimepb.FeedEntity) {
			assert.Nil(t, ent.Alert)
			assert.Nil(t, ent.IsDeleted)
			assert.Nil(t, ent.TripModifications)
			assert.Nil(t, ent.Shape)
			assert.Nil(t, ent.Vehicle)
			assert.Nil(t, ent.Stop)
			assert.NotNil(t, ent.Id)
			assert.NotNil(t, ent.TripUpdate)

			assert.Nil(t, ent.TripUpdate.Delay)
			assert.Nil(t, ent.TripUpdate.TripProperties)
			assert.NotNil(t, ent.TripUpdate.Vehicle)
			assert.NotNil(t, ent.TripUpdate.StopTimeUpdate)
			assert.NotNil(t, ent.TripUpdate.Trip)

			trip := ent.TripUpdate.Trip

			assert.Nil(t, trip.ModifiedTrip)
			assert.Nil(t, trip.StartDate)
			assert.Nil(t, trip.StartTime)
			assert.NotNil(t, trip.TripId)

			for _, stu := range ent.TripUpdate.StopTimeUpdate {
				assert.Nil(t, stu.DepartureOccupancyStatus)
				assert.Nil(t, stu.StopTimeProperties)
				assert.NotNil(t, stu.StopId)
				if stu.Arrival != nil {
					assert.Nil(t, stu.Arrival.Uncertainty)
					assert.Nil(t, stu.Arrival.ScheduledTime)
					assert.NotNil(t, stu.Arrival.Time)
				}
				if stu.Departure != nil {
					assert.Nil(t, stu.Departure.Uncertainty)
					assert.Nil(t, stu.Departure.ScheduledTime)
					assert.NotNil(t, stu.Departure.Time)
				}
			}

		}

		for _, ent := range feed.Entity {
			assertCommon(ent)
			switch source {
			case "A":
				assert.Nil(t, ent.TripUpdate.Timestamp)
				trip := ent.TripUpdate.Trip
				assert.Nil(t, trip.DirectionId)
				assert.Nil(t, trip.RouteId)
				assert.Nil(t, trip.ScheduleRelationship)

				for _, stu := range ent.TripUpdate.StopTimeUpdate {
					assert.Nil(t, stu.ScheduleRelationship)
					assert.Nil(t, stu.StopSequence)
					if stu.Arrival != nil {
						assert.Nil(t, stu.Arrival.Delay)
					}
					if stu.Departure != nil {
						assert.Nil(t, stu.Departure.Delay)
					}
				}
			case "T":
				assert.Nil(t, ent.TripUpdate.Timestamp)
				trip := ent.TripUpdate.Trip
				assert.Nil(t, trip.DirectionId)
				assert.Nil(t, trip.RouteId)
				assert.Nil(t, trip.ScheduleRelationship)

				for _, stu := range ent.TripUpdate.StopTimeUpdate {
					assert.NotNil(t, stu.ScheduleRelationship)
					assert.Equal(t, realtimepb.TripUpdate_StopTimeUpdate_SCHEDULED, *stu.ScheduleRelationship)
					assert.Nil(t, stu.StopSequence)
					if stu.Arrival != nil {
						assert.Nil(t, stu.Arrival.Delay)
					}
					if stu.Departure != nil {
						assert.Nil(t, stu.Departure.Delay)
					}
				}
			case "M":
				assert.NotNil(t, ent.TripUpdate.Timestamp)
				trip := ent.TripUpdate.Trip
				assert.NotNil(t, trip.DirectionId)
				assert.NotNil(t, trip.RouteId)
				assert.NotNil(t, trip.ScheduleRelationship)
				assert.Equal(t, realtimepb.TripDescriptor_SCHEDULED, *trip.ScheduleRelationship)

				for _, stu := range ent.TripUpdate.StopTimeUpdate {
					assert.NotNil(t, stu.ScheduleRelationship)
					assert.Equal(t, realtimepb.TripUpdate_StopTimeUpdate_SCHEDULED, *stu.ScheduleRelationship)
					assert.NotNil(t, stu.StopSequence)
					if stu.Arrival != nil {
						assert.NotNil(t, stu.Arrival.Delay)
					}
					if stu.Departure != nil {
						assert.NotNil(t, stu.Departure.Delay)
					}
				}
			}
		}
	}
}

func TestParseUpdates(t *testing.T) {
	b, err := os.ReadFile("testdata/TripUpdates_A.pb")
	require.NoError(t, err)
	feed, err := ParseFeed(b)
	require.NoError(t, err)
	updates, err := Bus.ParseStopUpdates(feed)
	require.NoError(t, err)
	_ = updates
}
