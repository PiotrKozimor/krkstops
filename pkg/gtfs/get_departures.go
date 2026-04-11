package gtfs

import (
	"slices"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	dayMinutes = 24 * 60
)

func (d *Departures) Get(stopName string, since time.Time, forMinutes int, filters ...DirectedRoute) ([]Departure, []RouteHeadsign) {
	stops, ok := d.Stops[stopName]
	if !ok {
		return nil, nil
	}

	since = since.In(Location)
	sinceMinutes := since.Hour()*60 + since.Minute()
	untilMinutes := sinceMinutes + forMinutes

	var departures []Departure
	serviceId := d.serviceId(since)

	forService := func(since, until int, serviceId uint32) []Departure {
		var departures []Departure
		for _, stopId := range stops {
			allDepartures, ok := d.lookup[departureKey{
				serviceId: serviceId,
				stopId:    stopId,
			}]
			if !ok {
				continue
			}

			sinceIndex, _ := slices.BinarySearchFunc(allDepartures, since, func(a departure, b int) int {
				return int(a.PlannedMinutesInDay) - b
			})
			untilIndex, _ := slices.BinarySearchFunc(allDepartures, until, func(a departure, b int) int {
				return int(a.PlannedMinutesInDay) - b
			})
			if untilIndex < sinceIndex {
				untilIndex = len(allDepartures)
			}
			selected := allDepartures[sinceIndex:untilIndex]

			d.updatesMu.RLock()
			for _, dep := range selected {
				departure := Departure{
					departure: dep,
				}
				if seconds, ok := d.updates[stopUpdateKey{
					stopId: stopId,
					tripId: departure.TripId,
				}]; ok {
					departure.UpdatedSecondsInDay = seconds
				}
				departures = append(departures, departure)
			}
			d.updatesMu.RUnlock()
		}
		return departures
	}

	departures = append(departures, forService(sinceMinutes, untilMinutes, serviceId)...)
	if since.Hour() < 5 {
		serviceId := d.serviceId(since.Add(-time.Hour * 24))
		departures = append(departures, forService(sinceMinutes+dayMinutes, untilMinutes+dayMinutes, serviceId)...)
	}

	slices.SortFunc(departures, func(a, b Departure) int {
		return int(a.PlannedMinutesInDay) - int(b.PlannedMinutesInDay)
	})

	routeToHeadsign := make(map[DirectedRoute]string, len(filters))

	for _, f := range filters {
		routeToHeadsign[f] = ""
	}

	if len(filters) > 0 {
		departures = slices.DeleteFunc(departures, func(d Departure) bool {
			_, ok := routeToHeadsign[DirectedRoute{
				RouteName:   d.RouteName,
				DirectionId: d.DirectionId,
			}]
			return !ok
		})
	}

	for _, dep := range departures {
		key := DirectedRoute{
			RouteName:   dep.RouteName,
			DirectionId: dep.DirectionId,
		}
		if _, ok := routeToHeadsign[key]; !ok {
			if trip, ok := d.trips[dep.TripId]; ok {
				routeToHeadsign[key] = trip.Headsign
			} else {
				gtfsErrors.With(prometheus.Labels{"err": "no_trip"}).Inc()
			}
		}
	}

	headsignToRoute := map[string][]DirectedRoute{}

	for k, v := range routeToHeadsign {
		headsignToRoute[v] = append(headsignToRoute[v], k)
	}
	headsigns := []RouteHeadsign{}
	for k, v := range headsignToRoute {
		headsigns = append(headsigns, RouteHeadsign{
			Headsign: k,
			Routes:   v,
		})
	}
	return departures, headsigns
}

func (d *Departures) serviceId(since time.Time) uint32 {
	weekday := since.Weekday()
	serviceId := uint32(0)
	for _, svc := range d.services {
		if slices.Contains(svc.Days, weekday) {
			serviceId = svc.Id
			break
		}
	}
	for _, exc := range d.serviceExceptions {
		if since.After(exc.StartsAt) && exc.Added {
			serviceId = exc.ServiceId
		}
	}
	return serviceId
}
