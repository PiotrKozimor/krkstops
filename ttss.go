package krkstops

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/pkg/gtfs"
)

func (s *KrkStopsServer) GetDepartures(ctx context.Context, req *pb.GetDeparturesRequest) (*pb.GetDeparturesResponse, error) {
	since := time.Now().Add(-time.Minute * 30)

	var resp pb.GetDeparturesResponse

	var filters []gtfs.DirectedRoute
	for _, f := range req.Filters {
		filters = append(filters, gtfs.DirectedRoute{
			RouteName:   f.RouteName,
			DirectionId: f.DirectionId,
		})
	}
	s.departuresMu.RLock()
	defer s.departuresMu.RUnlock()
	for _, dep := range s.departures {
		departures, headsigns := dep.Get(req.StopName, since, 150, filters...)
		for _, d := range departures {
			resp.Departures = append(resp.Departures, &pb.Departure3{
				PlannedMinutesInDay: d.PlannedMinutesInDay,
				DirectionId:         d.DirectionId,
				RouteName:           d.RouteName,
				UpdatedSecondsInDay: d.UpdatedSecondsInDay,
				TripId:              d.TripId,
			})
		}
		for _, h := range headsigns {

			headsign := &pb.RouteHeadsign{
				Headsign: h.Headsign,
			}

			routes := make([]*pb.DirectedRoute, len(h.Routes))
			for i, r := range h.Routes {
				routes[i] = &pb.DirectedRoute{
					RouteName:   r.RouteName,
					DirectionId: r.DirectionId,
				}
			}

			i, found := slices.BinarySearchFunc(resp.Headsigns, headsign, func(a, b *pb.RouteHeadsign) int {
				return strings.Compare(a.Headsign, h.Headsign)
			})
			if !found {
				headsign.Routes = routes
				resp.Headsigns = slices.Insert(resp.Headsigns, i, headsign)
			} else {
				resp.Headsigns[i].Routes = append(resp.Headsigns[i].Routes, routes...)
			}
		}
	}
	return &resp, nil
}

func (s *KrkStopsServer) SearchStops(ctx context.Context, req *pb.SearchStopsRequest) (*pb.SearchStopsResponse, error) {
	s.searchMu.RLock()
	items := s.search.Search(req.Query, 10)
	s.searchMu.RUnlock()
	names := make([]string, len(items))
	for i := range items {
		names[i] = items[i].name
	}
	return &pb.SearchStopsResponse{
		Stops: names,
	}, nil
}
