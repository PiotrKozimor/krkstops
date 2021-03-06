package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/PiotrKozimor/krkstops/airly"
	"github.com/PiotrKozimor/krkstops/cache"
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/stops"
	"github.com/PiotrKozimor/krkstops/ttss"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v7"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

type krkStopsServer struct {
	pb.UnimplementedKrkStopsServer
	c stops.Clients
}

func (s *krkStopsServer) GetAirly(ctx context.Context, installation *pb.Installation) (*pb.Airly, error) {
	var a *pb.Airly
	var err error
	isCached, err := cache.IsAirlyCached(s.c.Redis, installation)
	if err != nil {
		log.Println(err)
		isCached = false
	}
	if !isCached {
		a, err = airly.GetAirly(installation)
		if err != nil {
			return a, err
		}
		go cache.CacheAirly(s.c.Redis, a, installation)
	} else {
		a, err = cache.GetCachedAirly(s.c.Redis, installation)
	}
	return a, err

	// return airly, err
}

func (s *krkStopsServer) GetDepartures(stop *pb.Stop, stream pb.KrkStops_GetDeparturesServer) error {

	var deps []pb.Departure
	isCached, err := cache.IsDepartureCached(s.c.Redis, stop)
	if err != nil {
		log.Println(err)
		isCached = false
	}
	if !isCached {
		deps, err = ttss.GetStopDepartures(stop)
		if err != nil {
			return err
		}
		go cache.CacheDepartures(s.c.Redis, deps, stop)
	} else {
		deps, err = cache.GetCachedDepartures(s.c.Redis, stop)
		ttl, err := s.c.Redis.TTL(cache.DepsPrefix + stop.ShortName).Result()
		livedFor := int32(cache.DepsExpire.Seconds() - ttl.Seconds())
		for index := range deps {
			if deps[index].RelativeTime != 0 {
				deps[index].RelativeTime -= livedFor
			}
		}
		if err != nil {
			return err
		}
	}

	for _, dep := range deps {
		if err := stream.Send(&dep); err != nil {
			return err
		}
	}
	return nil
}

func (s *krkStopsServer) SearchStops(search *pb.StopSearch, stream pb.KrkStops_SearchStopsServer) error {
	stops, err := s.c.RedisAutocompleter.SuggestOpts(
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
		name, err := s.c.Redis.Get(stop.Payload).Result()
		if err != nil {
			return err
		}
		if err := stream.Send(&pb.Stop{Name: name, ShortName: stop.Payload}); err != nil {
			return err
		}
	}
	return nil
}

func (s *krkStopsServer) FindNearestAirlyInstallation(ctx context.Context, location *pb.InstallationLocation) (*pb.Installation, error) {
	inst, err := airly.FindAirlyInstallation(location)
	return inst, err
}

func (s *krkStopsServer) GetAirlyInstallation(ctx context.Context, installation *pb.Installation) (*pb.Installation, error) {
	inst, err := airly.GetAirlyInstallation(installation)
	return inst, err
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	server := krkStopsServer{
		c: stops.Clients{
			RedisAutocompleter: redisearch.NewAutocompleter(os.Getenv("REDIS_ENDPOINT"), "search-stops"),
			Redis:              redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ENDPOINT")})}}
	pb.RegisterKrkStopsServer(grpcServer, &server)
	grpc_prometheus.Register(grpcServer)
	handler := promhttp.Handler()
	http.Handle("/metrics", handler)
	go func() {
		err := http.ListenAndServe(":8040", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Fatal(grpcServer.Serve(lis))
}
