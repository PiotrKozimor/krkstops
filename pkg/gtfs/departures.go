package gtfs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"slices"
	"sync"
)

type departure struct {
	RouteName           uint32
	TripId              uint32
	DirectionId         uint32
	PlannedMinutesInDay uint32
}

func (d departure) String() string {
	minutes := d.PlannedMinutesInDay % 60
	hours := (d.PlannedMinutesInDay - minutes) / 60
	return fmt.Sprintf("%d %d %d %02d:%02d %d", d.RouteName, d.DirectionId, d.TripId, hours, minutes, d.PlannedMinutesInDay)
}

type Departure struct {
	departure
	UpdatedSecondsInDay uint32
}

func (d Departure) String() string {
	return fmt.Sprintf("%v %d:%02d:%02d", d.departure, d.UpdatedSecondsInDay/3600, d.UpdatedSecondsInDay/60%60, d.UpdatedSecondsInDay%60)

}

type DirectedRoute struct {
	RouteName   uint32
	DirectionId uint32
}

type RouteHeadsign struct {
	Routes   []DirectedRoute
	Headsign string
}

type departureKey struct {
	serviceId uint32
	stopId    uint32
}

type departureLookup map[departureKey][]departure

type Departures struct {
	*Unmarshaler
	Name string

	lookup            departureLookup
	trips             Trips
	routes            Routes
	Stops             Stops
	StopsScore        map[uint32]uint32
	services          []Service
	serviceExceptions []ServiceException
	updates           StopUpdates
	updatesMu         sync.RWMutex
}

func NewDepartures(u *Unmarshaler, name string) *Departures {
	return &Departures{
		Unmarshaler: u,
		Name:        name,
		lookup:      make(departureLookup, 100000),
		updates:     make(StopUpdates, 1000),
		StopsScore:  make(map[uint32]uint32, 1000),
	}
}

func (d *Departures) SetStopUpdates(updates StopUpdates) {
	d.updatesMu.Lock()
	defer d.updatesMu.Unlock()
	d.updates = updates
}

func (d *Departures) Init(retrieve func(file string) (*csv.Reader, error)) error {
	d.serviceIds = make(map[string]uint32, 20)
	d.tripIds = make(map[string]uint32, d.estimatedTripSize)
	d.stopIds = make(map[string]uint32, 1000)
	d.routeHeadSigns = make(map[string][]string, 100)

	exceptions, err1 := retrieve("calendar_dates.txt")
	service, err2 := retrieve("calendar.txt")
	routes, err3 := retrieve("routes.txt")
	if err := errors.Join(err1, err2, err3); err != nil {
		return fmt.Errorf("retrieve initial: %w", err)
	}

	d.services, err2 = d.UnmarshalServices(service)
	d.serviceExceptions, err1 = d.UnmarshalServiceExceptions(exceptions)
	d.routes, err3 = d.UnmarshalRoutes(routes)
	if err := errors.Join(err1, err2, err3); err != nil {
		return fmt.Errorf("unmarshal initial: %w", err)
	}

	r, err := retrieve("stops.txt")
	if err != nil {
		return fmt.Errorf("retrieve stops: %w", err)
	}
	d.Stops, err = d.UnmarshalStops(r)
	if err != nil {
		return fmt.Errorf("unmarshal stops: %w", err)
	}

	r, err = retrieve("trips.txt")
	if err != nil {
		return fmt.Errorf("retrieve trips: %w", err)
	}
	d.trips, err = d.UnmarshalTrips(r)
	if err != nil {
		return fmt.Errorf("unmarshal trips: %w", err)
	}

	r, err = retrieve("stop_times.txt")
	if err != nil {
		return fmt.Errorf("retrieve stop times: %w", err)
	}
	err = d.UnmarshalStopsTimes(r)
	if err != nil {
		return fmt.Errorf("unmarshal stop times: %w", err)
	}

	for key := range d.lookup {
		slices.SortFunc(d.lookup[key], func(a, b departure) int {
			return int(a.PlannedMinutesInDay) - int(b.PlannedMinutesInDay)
		})
	}
	return nil
}
