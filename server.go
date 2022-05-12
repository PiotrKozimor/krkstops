package krkstops

import (
	"context"
	"log"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/pkg/airly"
	"github.com/PiotrKozimor/krkstops/pkg/cache"
	"github.com/PiotrKozimor/krkstops/pkg/ttss"
	_ "github.com/kelseyhightower/envconfig"
)

var ENDPOINT = "krkstops-1.germanywestcentral.cloudapp.azure.com:8080"

type KrkStopsServer struct {
	pb.UnimplementedKrkStopsServer
	cache *cache.Cache
	Airly airly.Endpoint
	Ttss  []ttss.Endpointer
}

func NewServer(redisURI string) *KrkStopsServer {
	cache := cache.NewCache(redisURI, cache.SUG)
	return &KrkStopsServer{
		cache: cache,
	}
}

func (s *KrkStopsServer) GetAirly(ctx context.Context, installation *pb.Installation) (*pb.Airly, error) {
	var a *pb.Airly
	var err error
	conn := s.cache.Pool.Get()
	defer conn.Close()
	a, err = s.cache.GetAirly(installation, conn)
	if err != nil {
		a, err = s.Airly.GetAirly(installation)
		if err != nil {
			return a, err
		}
		err = s.cache.Airly(a, installation, cache.AirlyExpire, conn)
		if err != nil {
			log.Printf("cache airly: %v", err)
		}
	}
	return a, err
}

func (s *KrkStopsServer) GetDepartures2(ctx context.Context, stop *pb.Stop) (*pb.Departures, error) {
	conn := s.cache.Pool.Get()
	defer conn.Close()
	cachedDeps, err := s.cache.GetDepartures(ctx, stop, cache.DepsExpire, conn)
	if err != nil {
		endpoints, err := s.cache.GetEndpoints(ctx, stop.Id)
		if err != nil {
			return nil, err
		}
		deps := make([][]pb.Departure, len(endpoints))
		allDepsLen := 0
		for i, e := range endpoints {
			deps[i], err = s.Ttss[e].GetDepartures(uint(stop.Id))
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
		err = s.cache.Departures(&allDeps, stop, cache.DepsExpire, conn)
		if err != nil {
			log.Printf("cache departures: %v", err)
		}
		return &allDeps, nil
	} else {
		return cachedDeps, nil
	}
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
