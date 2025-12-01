package gtfs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Trip struct {
	Headsign    string
	DirectionId uint32
	ServiceId   uint32
	RouteId     uint32
}

type Trips map[uint32]Trip

func (u *Unmarshaler) UnmarshalTrips(r *csv.Reader) (Trips, error) {
	trips := make(Trips, u.estimatedTripSize)
	header, _ := r.Read()
	header[0] = strings.TrimPrefix(header[0], "\ufeff")

	idIndex := slices.Index(header, "trip_id")
	if idIndex < 0 {
		return nil, errors.New("trip_id index not found")
	}
	serviceIdIndex := slices.Index(header, "service_id")
	if idIndex < 0 {
		return nil, errors.New("trip_id index not found")
	}
	routeIdIndex := slices.Index(header, "route_id")
	if idIndex < 0 {
		return nil, errors.New("route_id index not found")
	}

	err := u.iterate(r, func(record []string) (err error) {
		trip := Trip{}
		trip.Headsign = record[3]
		id, err := u.reduceTripId(record[idIndex])
		if err != nil {
			return fmt.Errorf("reduce trip id: %w", err)
		}

		trip.ServiceId, err = u.reduceServiceId(record[serviceIdIndex])
		if err != nil {
			return fmt.Errorf("reduce service id: %w", err)
		}

		trip.RouteId, err = u.reduceRouteId(record[routeIdIndex])
		if err != nil {
			return fmt.Errorf("reduce route id: %w", err)
		}

		dirId, err := strconv.Atoi(record[5])
		if dirId != 0 && dirId != 1 {
			return fmt.Errorf("unexpected direction id: %d", dirId)
		}
		if err != nil {
			return fmt.Errorf("parse direction id: %w", err)
		}
		trip.DirectionId = uint32(dirId)
		trips[id] = trip
		return nil
	})
	return trips, err

}
