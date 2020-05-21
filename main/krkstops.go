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
		if err := stream.Send(&pb.Stop{Name: stop.Term, ShortName: stop.Payload}); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	server := krkStopsServer{
		app: krkstops.App{
			HTTPClient:         &http.Client{},
			RedisAutocompleter: redisearch.NewAutocompleter("localhost:6379", "search-stops"),
			RedisClient:        redis.NewClient(&redis.Options{Addr: "localhost:6379"})}}
	pb.RegisterKrkStopsServer(grpcServer, &server)
	grpcServer.Serve(lis)
}
