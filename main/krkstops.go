package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/PiotrKozimor/krk-stops-backend-golang/krkstops"
	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v7"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

type krkStopsServer struct {
	pb.UnimplementedKrkStopsServer
	app krkstops.App
}

func (s *krkStopsServer) GetAirly(ctx context.Context, installation *pb.Installation) (*pb.Airly, error) {
	airly, err := s.app.GetAirly(installation)
	return &airly, err
}

func (s *krkStopsServer) GetDepartures(stop *pb.Stop, stream pb.KrkStops_GetDeparturesServer) error {
	deps, err := s.app.GetStopDepartures(stop)
	if err != nil {
		return err
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
			RedisAutocompleter: redisearch.NewAutocompleter("localhost:6379", "search-stops"),
			RedisClient:        redis.NewClient(&redis.Options{Addr: "localhost:6379"})}}
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
