package krkstops

import (
	"context"
	"log"
	"strconv"

	"github.com/PiotrKozimor/krkstops/airly"
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/ttss"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/sirupsen/logrus"
)

var ENDPOINT = "krkstops.germanywestcentral.cloudapp.azure.com:8080"

type KrkStopsServer struct {
	pb.UnimplementedKrkStopsServer
	C     Clients
	Airly airly.Endpoint
	Ttss  []ttss.Endpointer
}

func (s *KrkStopsServer) GetAirly(ctx context.Context, installation *pb.Installation) (*pb.Airly, error) {
	var a *pb.Airly
	var err error
	isCached, err := isAirlyCached(s.C.Redis, installation)
	if err != nil {
		log.Println(err)
		isCached = false
	}
	if !isCached {
		a, err = s.Airly.GetAirly(installation)
		if err != nil {
			return a, err
		}
		go cacheAirly(s.C.Redis, a, installation)
	} else {
		a, err = getCachedAirly(s.C.Redis, installation)
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
	stops, err := s.C.RedisAutocompleter.SuggestOpts(
		search.Query, redisearch.SuggestOptions{
			Num:          10,
			Fuzzy:        true,
			WithPayloads: true,
			WithScores:   false,
		})
	if err != nil {
		return err
	}
	for _, stop := range stops {
		name, err := s.C.Redis.Get(ctx, stop.Payload).Result()
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(stop.Payload)
		if err != nil {
			logrus.Errorf("failed to parse %s to int", stop.Payload)
		} else if err := stream.Send(&pb.Stop{Name: name, Id: uint32(id)}); err != nil {
			return err
		}
	}
	return nil
}

func (s *KrkStopsServer) SearchStops2(ctx context.Context, search *pb.StopSearch) (*pb.Stops, error) {
	stops, err := s.C.RedisAutocompleter.SuggestOpts(
		search.Query, redisearch.SuggestOptions{
			Num:          10,
			Fuzzy:        true,
			WithPayloads: true,
			WithScores:   false,
		})
	if err != nil {
		return nil, err
	}
	stopsP := make([]*pb.Stop, len(stops))
	for i, stop := range stops {
		name, err := s.C.Redis.Get(ctx, stop.Payload).Result()
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(stop.Payload)
		if err != nil {
			logrus.Errorf("failed to parse %s to int", stop.Payload)
		} else {
			stopsP[i] = &pb.Stop{Name: name, Id: uint32(id)}
		}
	}
	return &pb.Stops{
		Stops: stopsP,
	}, nil
}

func (s *KrkStopsServer) FindNearestAirlyInstallation(ctx context.Context, location *pb.InstallationLocation) (*pb.Installation, error) {
	inst, err := s.Airly.NearestInstallation(location)
	return inst, err
}

func (s *KrkStopsServer) GetAirlyInstallation(ctx context.Context, installation *pb.Installation) (*pb.Installation, error) {
	inst, err := s.Airly.GetInstallation(uint(installation.Id))
	return inst, err
}
