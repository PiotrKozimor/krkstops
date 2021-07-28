package krkstops

import (
	"context"

	"github.com/PiotrKozimor/krkstops/airly"
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/ttss"
)

var ENDPOINT = "krkstops.germanywestcentral.cloudapp.azure.com:8080"

type KrkStopsServer struct {
	pb.UnimplementedKrkStopsServer
	cache *Cache
	Airly airly.Endpoint
	Ttss  []ttss.Endpointer
}

func NewServer(redisURI string) (*KrkStopsServer, error) {
	cache, err := NewCache(redisURI, SUG)
	if err != nil {
		return nil, err
	}
	server := KrkStopsServer{
		cache: cache,
		Airly: airly.Api,
		Ttss:  ttss.KrkStopsEndpoints,
	}
	return &server, nil
}

func (s *KrkStopsServer) GetAirly(ctx context.Context, installation *pb.Installation) (*pb.Airly, error) {
	var a *pb.Airly
	var err error
	a, err = s.cache.getAirly(installation)
	if err != nil {
		a, err = s.Airly.GetAirly(installation)
		if err != nil {
			return a, err
		}
		go s.cache.airly(a, installation)
	}
	return a, err
}

func (s *KrkStopsServer) GetDepartures(stop *pb.Stop, stream pb.KrkStops_GetDeparturesServer) error {
	deps, err := s.cache.getDepartures(context.Background(), stop)
	if err != nil {
		depsC, errC := ttss.GetDepartures(s.Ttss, uint(stop.Id))
		for d := range depsC {
			for i := range d {
				err := stream.Send(&d[i])
				if err != nil {
					return err
				}
				deps.Departures = append(deps.Departures, &d[i])
			}
		}
		for err := range errC {
			return err
		}
		go s.cache.departures(deps, stop)
	} else {
		for _, dep := range deps.Departures {
			if err := stream.Send(dep); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *KrkStopsServer) GetDepartures2(ctx context.Context, stop *pb.Stop) (*pb.Departures, error) {
	cachedDeps, err := s.cache.getDepartures(ctx, stop)
	if err != nil {
		endpoints, err := s.cache.getEndpoints(ctx, stop.Id)
		if err != nil {
			return nil, err
		}
		deps := make([][]pb.Departure, len(endpoints))
		allDepsLen := 0
		for i, e := range endpoints {
			deps[i], err = e.GetDepartures(uint(stop.Id))
			if err != nil {
				return nil, err
			}
			allDepsLen += len(deps[i])
		}
		allDeps := pb.Departures{
			Departures: make([]*pb.Departure, allDepsLen),
		}
		k := 0
		for i := range deps {
			for j := range deps[i] {
				allDeps.Departures[k] = &deps[i][j]
				k++
			}
		}
		go s.cache.departures(&allDeps, stop)
		return &allDeps, nil
	} else {
		return cachedDeps, nil
	}
}

func (s *KrkStopsServer) SearchStops(search *pb.StopSearch, stream pb.KrkStops_SearchStopsServer) error {
	stops, err := s.cache.Search(context.Background(), search.Query)
	if err != nil {
		return err
	}
	for _, stop := range stops {
		if err := stream.Send(stop); err != nil {
			return err
		}
	}
	return nil
}

func (s *KrkStopsServer) SearchStops2(ctx context.Context, search *pb.StopSearch) (*pb.Stops, error) {
	stops, err := s.cache.Search(ctx, search.Query)
	return &pb.Stops{
		Stops: stops,
	}, err
}

func (s *KrkStopsServer) FindNearestAirlyInstallation(ctx context.Context, location *pb.InstallationLocation) (*pb.Installation, error) {
	inst, err := s.Airly.NearestInstallation(location)
	return inst, err
}

func (s *KrkStopsServer) GetAirlyInstallation(ctx context.Context, installation *pb.Installation) (*pb.Installation, error) {
	inst, err := s.Airly.GetInstallation(uint(installation.Id))
	return inst, err
}
