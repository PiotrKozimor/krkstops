package gtfs

import (
	"log"
	"strconv"
	"strings"

	"github.com/artonge/go-gtfs"
	"github.com/sirupsen/logrus"
)

type Departure struct {
	StopId      int
	PatternText string
	Direction   string
}

type DeparturesDirections map[Departure][]int

func Parse(data *gtfs.GTFS) (DeparturesDirections, error) {
	var routeIdToPatternTextLookup map[string]string
	var tripIdToDirectionsLookup map[string]string
	var tripIdToRouteIdLookup map[string]string
	dirs := make(DeparturesDirections, 500)
	routeIdToPatternTextLookup = make(map[string]string, len(data.Routes))
	tripIdToDirectionsLookup = make(map[string]string, len(data.Trips))
	tripIdToRouteIdLookup = make(map[string]string, len(data.Trips))
	for i := range data.Routes {
		routeIdToPatternTextLookup[data.Routes[i].ID] = data.Routes[i].ShortName
	}
	// log.Printf("routeIdToPatternTextLookup: %#+v\n", routeIdToPatternTextLookup)
	for i := range data.Trips {
		tripIdToDirectionsLookup[data.Trips[i].ID] = data.Trips[i].Headsign
	}
	// log.Printf("tripIdToDirectionsLookup: %#+v\n", tripIdToDirectionsLookup)
	for i := range data.Trips {
		tripIdToRouteIdLookup[data.Trips[i].ID] = data.Trips[i].RouteID
	}
	// log.Printf("tripIdToRouteIdLookup: %#+v\n", tripIdToRouteIdLookup)
	for i := range data.StopsTimes {
		ttssIdPostfixed := strings.Split(data.StopsTimes[i].StopID, "_")[2]
		ttssIdStr := ttssIdPostfixed[0 : len(ttssIdPostfixed)-2]
		parsedStopDirection, err := strconv.Atoi(ttssIdPostfixed[len(ttssIdPostfixed)-2:])
		if err != nil {
			return nil, err
		}
		ttssId, err := strconv.Atoi(ttssIdStr)
		if err != nil {
			return nil, err
		}
		tripId := data.StopsTimes[i].TripID
		dep := Departure{
			StopId:      ttssId,
			PatternText: routeIdToPatternTextLookup[tripIdToRouteIdLookup[tripId]],
			Direction:   tripIdToDirectionsLookup[tripId],
		}
		if dirs[dep] == nil {
			// TODO When bus or tram departs from stop with given name, here we suppose that it can to it maximum three times, e.g.
			// 194 (heading to Krowodrza Górka) on Rondo Grunwaldzkie stop.
			dirs[dep] = make([]int, 3)
			dirs[dep][0] = parsedStopDirection
			logrus.Debugf("INITSLICE: group: %v dep: %v, trip_id: %s", dirs[dep], dep, tripId)
		} else {
			saved := false
			for i := range dirs[dep] {
				if dirs[dep][i] == parsedStopDirection {
					saved = true
					logrus.Debugf("OK: group: %d dep: %v", dirs[dep], dep)
					break
				} else if dirs[dep][i] == 0 {
					dirs[dep][i] = parsedStopDirection
					saved = true
					logrus.Infof("INIT: group: %v dep: %v, trip_id: %s", dirs[dep], dep, tripId)
					break
				}
			}
			if !saved {
				logrus.Warnf("CONFLICT: group: %d parsed: %d dep: %v, trip_id: %s", dirs[dep], parsedStopDirection, dep, tripId)
			}
		}
	}
	log.Print(len(dirs))
	return dirs, nil
}

func ParseMany(data ...*gtfs.GTFS) (DeparturesDirections, error) {
	grouping, err := Parse(data[0])
	if err != nil {
		return nil, err
	}
	if len(data) > 1 {
		for i := 1; i < len(data); i++ {
			tmpGrouping, err := Parse(data[i])
			if err != nil {
				return nil, err
			}
			for k := range tmpGrouping {
				grouping[k] = tmpGrouping[k]
			}
		}
	}
	return grouping, nil
}
