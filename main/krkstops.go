package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/PiotrKozimor/krk-stops-backend-golang/cache"
	"github.com/PiotrKozimor/krk-stops-backend-golang/krkstops"
	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v7"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

type krkStopsServer struct {
	pb.UnimplementedKrkStopsServer
	app krkstops.App
}

func (s *krkStopsServer) GetAirly(ctx context.Context, installation *pb.Installation) (*pb.Airly, error) {
	var airly pb.Airly
	var err error
	isCached, err := cache.IsAirlyCached(s.app.RedisClient, installation)
	if err != nil {
		log.Println(err)
		isCached = false
	}
	if !isCached {
		airly, err = s.app.GetAirly(installation)
		if err != nil {
			return &airly, err
		}
		go cache.CacheAirly(s.app.RedisClient, &airly, installation)
	} else {
		airly, err = cache.GetCachedAirly(s.app.RedisClient, installation)
	}
	return &airly, err

	// return airly, err
}

func (s *krkStopsServer) GetDepartures(stop *pb.Stop, stream pb.KrkStops_GetDeparturesServer) error {

	var deps []pb.Departure
	isCached, err := cache.IsDepartureCached(s.app.RedisClient, stop)
	if err != nil {
		log.Println(err)
		isCached = false
	}
	if !isCached {
		deps, err = s.app.GetStopDepartures(stop)
		if err != nil {
			return err
		}
		go cache.CacheDepartures(s.app.RedisClient, &deps, stop)
	} else {
		deps, err = cache.GetCachedDepartures(s.app.RedisClient, stop)
		ttl, err := s.app.RedisClient.TTL(cache.DepsPrefix + stop.ShortName).Result()
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
	stops, err := s.app.RedisAutocompleter.SuggestOpts(
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
		name, err := s.app.RedisClient.Get(stop.Payload).Result()
		if err != nil {
			return err
		}
		if err := stream.Send(&pb.Stop{Name: name, ShortName: stop.Payload}); err != nil {
			return err
		}
	}
	return nil
}

func (s *krkStopsServer) GetAirlyInstallation(ctx context.Context, location *pb.InstallationLocation) (*pb.Installation, error) {
	inst, err := s.app.GetAirlyInstallation(location)
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
		app: krkstops.App{
			HTTPClient:         &http.Client{},
			RedisAutocompleter: redisearch.NewAutocompleter(os.Getenv("REDIS_ENDPOINT"), "search-stops"),
			RedisClient:        redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ENDPOINT")})}}
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
