package krkstops

import (
	"context"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/pkg/search"
	"github.com/PiotrKozimor/krkstops/pkg/ttss"
	"golang.org/x/exp/slices"
)

type typedDeparture struct {
	ttss.Departure
	transit pb.Transit
}

type typedClient struct {
	ttssClient
	transit pb.Transit
}

func (s *KrkStopsServer) GetDepartures2(ctx context.Context, req *pb.GetDepartures2Request) (*pb.GetDepartures2Response, error) {
	cachedDeps, cachedAt, ok := s.depsCache.get(uint(req.Id))
	if !ok {
		stop := s.searchCli.Get(uint(req.Id))
		clients := []typedClient{}
		if stop.Bus {
			clients = append(clients, typedClient{s.ttssBusCli, pb.Transit_BUS})
		}
		if stop.Tram {
			clients = append(clients, typedClient{s.ttssTramCli, pb.Transit_TRAM})
		}
		for _, typedCli := range clients {
			deps, err := typedCli.GetDepartures(stop.Id)
			if err != nil {
				return nil, err
			}
			if cachedDeps == nil {
				cachedDeps = make([]typedDeparture, 0, len(deps))
			}
			for i := range deps {
				cachedDeps = append(cachedDeps, typedDeparture{
					Departure: deps[i],
					transit:   typedCli.transit,
				})
			}
			slices.SortFunc(cachedDeps, func(a, b typedDeparture) int {
				return int(a.RelativeTime - b.RelativeTime)
			})
		}
		s.depsCache.set(stop.Id, cachedDeps)
	}
	return &pb.GetDepartures2Response{
		Departures: protoDepartures(cachedDeps, cachedAt),
	}, nil
}

func (s *KrkStopsServer) SearchStops2(ctx context.Context, req *pb.SearchStops2Request) (*pb.SearchStops2Response, error) {
	stops := s.searchCli.Search(req.Query, 10)
	return &pb.SearchStops2Response{
		Stops: protoStops(stops),
	}, nil
}

func protoStops(stops []search.Stop) []*pb.Stop {
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
