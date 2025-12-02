package gtfs

import (
	"slices"
	"time"
)

func (d *Departures) Get(stopName string, since, until time.Time, filters ...DirectedRoute) ([]Departure, []RouteHeadsign) {
	stops, ok := d.Stops[stopName]
	if !ok {
		return nil, nil
	}

	since = since.In(location)
	until = until.In(location)
	sinceMinutes := since.Hour()*60 + since.Minute()
	untilMinutes := until.Hour()*60 + until.Minute()

	var departures []Departure
	serviceId := d.serviceId(since)

	for _, stopId := range stops {
		allDepartures, ok := d.lookup[departureKey{
			serviceId: serviceId,
			stopId:    stopId,
		}]
		if !ok {
			continue
		}

		sinceIndex, _ := slices.BinarySearchFunc(allDepartures, sinceMinutes, func(a departure, b int) int {
			return int(a.PlannedMinutesInDay) - b
		})
		untilIndex, _ := slices.BinarySearchFunc(allDepartures, untilMinutes, func(a departure, b int) int {
			return int(a.PlannedMinutesInDay) - b
		})
		if untilIndex < sinceIndex {
			untilIndex = len(allDepartures)
		}
		selected := allDepartures[sinceIndex:untilIndex]

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
				// TODO count
			}
		} else {
			// TODO count
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
	if serviceId == 0 {
		for _, exc := range d.serviceExceptions {
			if since.After(exc.StartsAt) {
				serviceId = exc.ServiceId
			}
		}
	}
	return serviceId
}
