package krkstops

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/pkg/gtfs"
	"github.com/PiotrKozimor/krkstops/pkg/ttss"
	"github.com/PiotrKozimor/krkstops/pkg/ttssstops"
	"google.golang.org/grpc"
)

type typedDeparture struct {
	ttss.Departure
	transit pb.Transit
}

func (s *KrkStopsServer) GetDepartures2(ctx context.Context, req *pb.GetDepartures2Request) (*pb.GetDepartures2Response, error) {
	cachedDeps, cachedAt, ok := s.depsCache.get(uint(req.Id))
	if !ok {
		stop := s.searchCli.Get(uint(req.Id))
		cachedDeps = make([]typedDeparture, 0, 20)
		errs := make([]error, 0)
		m := sync.Mutex{}
		wg := sync.WaitGroup{}

		getDepartures := func(cli ttssClient, transit pb.Transit) {
			deps, err := cli.GetDepartures(stop.Id)
			m.Lock()
			if err != nil {
				errs = append(errs, err)
			} else {
				for _, dep := range deps {
					cachedDeps = append(cachedDeps, typedDeparture{
						Departure: dep,
						transit:   transit,
					})
				}

			}
			m.Unlock()
			wg.Done()
		}
		if stop.Bus {
			wg.Add(1)
			go getDepartures(s.ttssBusCli, pb.Transit_BUS)
		}
		if stop.Tram {
			wg.Add(1)
			go getDepartures(s.ttssTramCli, pb.Transit_TRAM)
		}

		wg.Wait()
		err := errors.Join(errs...)
		if err != nil {
			if len(cachedDeps) == 0 {
				return nil, err
			} else {
				slog.ErrorContext(ctx, "get departures", "error", err)
			}
		}

		slices.SortFunc(cachedDeps, func(a, b typedDeparture) int {
			return int(a.RelativeTime - b.RelativeTime)
		})
		s.depsCache.set(stop.Id, cachedDeps)
	}
	if len(cachedDeps) > 20 {
		grpc.SetSendCompressor(ctx, "gzip")
	}

	return &pb.GetDepartures2Response{
		Departures: protoDepartures(cachedDeps, cachedAt),
	}, nil
}

func (s *KrkStopsServer) GetDepartures3(ctx context.Context, req *pb.GetDepartures3Request) (*pb.GetDepartures3Response, error) {
	since := time.Now().Add(-time.Minute * 5)
	until := since.Add(time.Hour)

	var resp pb.GetDepartures3Response

	var filters []gtfs.DirectedRoute
	for _, f := range req.Filters {
		filters = append(filters, gtfs.DirectedRoute{
			RouteName:   f.RouteName,
			DirectionId: f.DirectionId,
		})
	}
	for _, dep := range s.departures {
		departures, headsigns := dep.Get(req.StopName, since, until, filters...)
		for _, d := range departures {
			resp.Departures = append(resp.Departures, &pb.Departure3{
				PlannedMinutesInDay: d.PlannedMinutesInDay,
				DirectionId:         d.DirectionId,
				RouteName:           d.RouteName,
				UpdatedSecondsInDay: d.UpdatedSecondsInDay,
				Transit:             pb.Transit2_BUS2, // TODO
			})
		}
		for _, h := range headsigns {
			headsign := &pb.RouteHeadsign{
				Headsign: h.Headsign,
			}
			for _, r := range h.Routes {
				headsign.Routes = append(headsign.Routes, &pb.DirectedRoute{
					RouteName:   r.RouteName,
					DirectionId: r.DirectionId,
				})
			}
			resp.Headsigns = append(resp.Headsigns, headsign)
		}
	}
	return &resp, nil
}

func (s *KrkStopsServer) SearchStops2(ctx context.Context, req *pb.SearchStops2Request) (*pb.SearchStops2Response, error) {
	stops := s.searchCli.Search(req.Query, 10)
	return &pb.SearchStops2Response{
		Stops: protoStops(stops),
	}, nil
}

func (s *KrkStopsServer) SearchStops3(ctx context.Context, req *pb.SearchStops3Request) (*pb.SearchStops3Response, error) {
	items := s.search.Search(req.Query, 10)
	names := make([]string, len(items))
	for i := range items {
		names[i] = items[i].name
	}
	return &pb.SearchStops3Response{
		Stops: names,
	}, nil
}

func protoStops(stops []ttssstops.Stop) []*pb.Stop {
	pStops := make([]*pb.Stop, len(stops))
	for i := range stops {
		pStops[i] = &pb.Stop{
			Name: stops[i].Name,
			Id:   uint32(stops[i].Id),
		}
	}
	return pStops
}

func protoDepartures(departures []typedDeparture, cachedAt *time.Time) []*pb.Departure {
	pDepartures := make([]*pb.Departure, len(departures))
	for i := range departures {
		pDepartures[i] = &pb.Departure{
			RelativeTime: departures[i].Departure.RelativeTime,
			PlannedTime:  departures[i].Departure.PlannedTime,
			PatternText:  departures[i].Departure.PatternText,
			Direction:    departures[i].Departure.Direction,
			Type:         departures[i].transit,
			Predicted:    departures[i].Departure.Predicted,
		}
		if cachedAt != nil {
			pDepartures[i].RelativeTime -= int32(time.Now().Sub(*cachedAt).Seconds())
		}
	}
	return pDepartures
}
