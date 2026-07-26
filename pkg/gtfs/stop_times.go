package gtfs

import (
	"encoding/csv"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
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
			return fmt.Errorf("route not found: %s", trip.RouteId)
		}

		directionId, found := slices.BinarySearch(d.routeHeadSigns[trip.RouteId], trip.Headsign)
		if !found {
			return fmt.Errorf("headsign not found: %s", trip.Headsign)
		}
		dep := departure{
			TripId:              tripId,
			DirectionId:         uint32(directionId),
			PlannedMinutesInDay: minutesInDay,
		}
		routeInt, err := strconv.Atoi(route)
		if err == nil {
			dep.RouteName = uint32(routeInt)
		} else if trimmed := strings.TrimFunc(route, func(r rune) bool { return !unicode.IsDigit(r) }); len(trimmed) > 0 {
			if routeInt2, err := strconv.Atoi(trimmed); err == nil {
				dep.RouteName = uint32(routeInt2)
			}
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
