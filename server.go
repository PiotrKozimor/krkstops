package krkstops

import (
	"context"
	"log"
	"strconv"

	"github.com/PiotrKozimor/krkstops/airly"
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/ttss"
	"github.com/go-redis/redis/v8"
)

var ENDPOINT = "krkstops.germanywestcentral.cloudapp.azure.com:8080"

type KrkStopsServer struct {
	pb.UnimplementedKrkStopsServer
	C     *Cache
	Airly airly.Endpoint
	Ttss  []ttss.Endpointer
}

func (s *KrkStopsServer) GetAirly(ctx context.Context, installation *pb.Installation) (*pb.Airly, error) {
	var a *pb.Airly
	var err error
	a, err = getCachedAirly(s.C.Redis, installation)
	if err == redis.Nil {
		a, err = s.Airly.GetAirly(installation)
		if err != nil {
			return a, err
		}
		go cacheAirly(s.C.Redis, a, installation)
	}
	return a, err

	// return airly, err
}

func (s *KrkStopsServer) GetDepartures(stop *pb.Stop, stream pb.KrkStops_GetDeparturesServer) error {

	var deps []pb.Departure
	isCached, err := isDepartureCached(s.C.Redis, stop)
	if err != nil {
		log.Println(err)
		isCached = false
	}
	if !isCached {
		depsC, errC := ttss.GetDepartures(s.Ttss, uint(stop.Id))
		for d := range depsC {
			for _, departure := range d {
				err := stream.Send(&departure)
				if err != nil {
					return err
				}
			}
			deps = append(deps, d...)
		}
		for err := range errC {
			return err
		}
		go cacheDepartures(s.C.Redis, deps, stop)
	} else {
		deps, err = getCachedDepartures(s.C.Redis, stop)
		ttl, err := s.C.Redis.TTL(ctx, depsPrefix+strconv.Itoa(int(stop.Id))).Result()
		livedFor := int32(depsExpire.Seconds() - ttl.Seconds())
		for index := range deps {
			if deps[index].RelativeTime != 0 {
				deps[index].RelativeTime -= livedFor
			}
		}
		if err != nil {
			return err
		}
		for _, dep := range deps {
			if err := stream.Send(&dep); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *KrkStopsServer) GetDepartures2(ctx context.Context, stop *pb.Stop) (*pb.Departures, error) {
	var deps []pb.Departure
	var returnedDepartures []*pb.Departure
	isCached, err := isDepartureCached(s.C.Redis, stop)
	if err != nil {
		log.Println(err)
		isCached = false
	}
	if !isCached {
		depsC, errC := ttss.GetDepartures(s.Ttss, uint(stop.Id))
		for d := range depsC {
			deps = append(deps, d...)
		}
		returnedDepartures = make([]*pb.Departure, len(deps))
		for i := range deps {
			returnedDepartures[i] = &deps[i]
		}
		for err := range errC {
			return nil, err
		}
		go cacheDepartures(s.C.Redis, deps, stop)
	} else {
		deps, err = getCachedDepartures(s.C.Redis, stop)
		if err != nil {
			return nil, err
		}
		ttl, err := s.C.Redis.TTL(ctx, depsPrefix+strconv.Itoa(int(stop.Id))).Result()
		if err != nil {
			return nil, err
		}
		livedFor := int32(depsExpire.Seconds() - ttl.Seconds())
		for index := range deps {
			if deps[index].RelativeTime != 0 {
				deps[index].RelativeTime -= livedFor
			}
		}
		returnedDepartures = make([]*pb.Departure, len(deps))
		for i := range deps {
			returnedDepartures[i] = &deps[i]
		}
	}
	return &pb.Departures{
		Departures: returnedDepartures,
	}, nil
}

func (s *KrkStopsServer) SearchStops(search *pb.StopSearch, stream pb.KrkStops_SearchStopsServer) error {
	stops, err := s.C.Search(search.Query)
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
	stops, err := s.C.Search(search.Query)
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
