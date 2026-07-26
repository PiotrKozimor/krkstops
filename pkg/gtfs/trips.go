package gtfs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"slices"
	"strings"
)

type Trip struct {
	Headsign  string
	RouteId   string
	ServiceId uint32
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
	headSign := slices.Index(header, "trip_headsign")
	if idIndex < 0 {
		return nil, errors.New("trip_headsign index not found")
	}

	err := u.iterateSorted(r, idIndex, func(i int, record []string) (err error) {
		trip := Trip{}
		trip.Headsign = record[headSign]
		id := uint32(i)
		u.tripIds[record[idIndex]] = id

		trip.ServiceId, err = u.serviceId(record[serviceIdIndex])
		if err != nil {
			return fmt.Errorf("reduce service id: %w", err)
		}

		trip.RouteId = record[routeIdIndex]

		headsigns := u.routeHeadSigns[trip.RouteId]
		i, found := slices.BinarySearch(headsigns, trip.Headsign)
		if !found {
			u.routeHeadSigns[trip.RouteId] = slices.Insert(headsigns, i, trip.Headsign)
		}
		trips[id] = trip
		return nil
	})
	return trips, err
}
