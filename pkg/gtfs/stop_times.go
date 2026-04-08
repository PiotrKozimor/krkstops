package gtfs

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

func (d *Departures) UnmarshalStopsTimes(r *csv.Reader) error {
	r.Read()
	return d.iterate(r, func(record []string) error {
		tripId, ok := d.Unmarshaler.tripIds[record[0]]
		if !ok {
			return fmt.Errorf("trip id not found: %s", record[0])
		}
		stopId, err := d.stopId(record[3])
		if err != nil {
			return fmt.Errorf("stop id: %w", err)
		}
		minutesInDay, err := parseTime(record[2])
		if err != nil {
			return fmt.Errorf("parse departure time: %w", err)
		}

		trip, ok := d.trips[tripId]
		if !ok {
			return fmt.Errorf("trip not found: %d", tripId)
		}

		route, ok := d.routes[trip.RouteId]
		if !ok {
			return fmt.Errorf("route not found: %d", trip.RouteId)
		}

		dep := departure{
			RouteName:           route.Name,
			TripId:              tripId,
			DirectionId:         trip.DirectionId,
			PlannedMinutesInDay: minutesInDay,
		}
		key := departureKey{
			serviceId: trip.ServiceId,
			stopId:    stopId,
		}
		d.StopsScore[stopId]++
		d.lookup[key] = append(d.lookup[key], dep)
		return nil
	})
}

func parseTime(inDay string) (uint32, error) {
	split := strings.Split(inDay, ":")
	if len(split) != 3 {
		return 0, fmt.Errorf("invalid time: %s", inDay)
	}
	hours, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hour: %w", err)
	}
	minute, err := strconv.Atoi(split[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minute: %w", err)
	}
	return uint32(hours*60 + minute), nil
}
