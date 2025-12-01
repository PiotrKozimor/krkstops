package gtfs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"slices"
	"sync"
	"time"
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
	return fmt.Sprintf("%d %d %02d:%02d", d.RouteName, d.DirectionId, hours, minutes)
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
	lookup            departureLookup
	trips             Trips
	routes            Routes
	Stops             Stops
	StopsScore        map[uint32]uint32
	services          []Service
	serviceExceptions []ServiceException
	updates           stopUpdates
	updatesMu         sync.RWMutex
	GetFeed           func() ([]byte, error)
}

func NewDepartures(u *Unmarshaler) *Departures {
	return &Departures{
		Unmarshaler: u,
		lookup:      make(departureLookup, 100000),
		updates:     make(stopUpdates, 1000),
		StopsScore:  make(map[uint32]uint32, 1000),
	}
}

func (d *Departures) RefreshUpdates() {
	t := time.NewTicker(time.Minute)
	for {
		err := d.fetchUpdates()
		if err != nil {
			log.Print("fetch updates: ", err)
		}
		<-t.C
	}
}

func (d *Departures) fetchUpdates() error {
	b, err := d.GetFeed()
	if err != nil {
		return fmt.Errorf("get feed: %w", err)
	}
	feed, err := ParseFeed(b)
	if err != nil {
		return fmt.Errorf("parse feed: %w", err)
	}
	updates, err := d.parseStopUpdates(feed)
	if err != nil {
		return fmt.Errorf("parse stop updates: %w", err)
	}
	d.updatesMu.RLock()
	defer d.updatesMu.RUnlock()
	d.updates = updates
	return nil
}

func (d *Departures) Init(retrieve func(file string) (*csv.Reader, error)) error {

	exceptions, err1 := retrieve("calendar_dates.txt")
	service, err2 := retrieve("calendar.txt")
	routes, err3 := retrieve("routes.txt")
	if err := errors.Join(err1, err2, err3); err != nil {
		return err
	}

	d.serviceExceptions, err1 = d.UnmarshalServiceExceptions(exceptions)
	d.services, err2 = d.UnmarshalServices(service)
	d.routes, err3 = d.UnmarshalRoutes(routes)
	if err := errors.Join(err1, err2, err3); err != nil {
		return err
	}

	r, err := retrieve("stops.txt")
	if err != nil {
		return err
	}
	d.Stops, err = d.UnmarshalStops(r)
	if err != nil {
		return err
	}

	r, err = retrieve("trips.txt")
	if err != nil {
		return err
	}
	d.trips, err = d.UnmarshalTrips(r)
	if err != nil {
		return err
	}

	r, err = retrieve("stop_times.txt")
	if err != nil {
		return err
	}
	err = d.UnmarshalStopsTimes(r)

	for key := range d.lookup {
		slices.SortFunc(d.lookup[key], func(a, b departure) int {
			return int(a.PlannedMinutesInDay) - int(b.PlannedMinutesInDay)
		})
	}
	return err
}
