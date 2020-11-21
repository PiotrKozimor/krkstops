package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"

	"github.com/PiotrKozimor/krkstops/pb"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type krkStopsServer struct {
	pb.UnimplementedKrkStopsServer
	returnErrors map[string]bool
}

var foobarError = errors.New("foo bar")

func (s *krkStopsServer) GetAirly(ctx context.Context, installation *pb.Installation) (*pb.Airly, error) {
	if s.returnErrors["airly"] {
		return nil, foobarError
	}
	return &pb.Airly{Caqi: 10, Humidity: 50, Temperature: 30.0, Color: "#999999"}, nil
}

func (s *krkStopsServer) GetDepartures(stop *pb.Stop, stream pb.KrkStops_GetDeparturesServer) error {
	if s.returnErrors["departures"] {
		return foobarError
	}
	if err := stream.Send(&pb.Departure{
		RelativeTime: 110,
		PlannedTime:  "14:48",
		Direction:    "Bronowice Małe",
		PatternText:  "8",
	}); err != nil {
		return err
	}
	return nil
}

func (s *krkStopsServer) SearchStops(search *pb.StopSearch, stream pb.KrkStops_SearchStopsServer) error {
	if s.returnErrors["search"] {
		return foobarError
	}
	if err := stream.Send(&pb.Stop{Name: "Norymberska", ShortName: "2688"}); err != nil {
		return err
	}
	return nil
}

func (s *krkStopsServer) FindNearestAirlyInstallation(ctx context.Context, location *pb.InstallationLocation) (*pb.Installation, error) {
	if s.returnErrors["searchInstallation"] {
		return nil, foobarError
	}
	return &pb.Installation{Id: 9914, Latitude: 50.013157, Longitude: 19.906305}, nil
}

func (s *krkStopsServer) GetAirlyInstallation(ctx context.Context, installation *pb.Installation) (*pb.Installation, error) {
	if s.returnErrors["searchInstallation"] {
		return nil, foobarError
	}
	return &pb.Installation{Id: 9914, Latitude: 50.013157, Longitude: 19.906305}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", ":10475")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	server := krkStopsServer{
		returnErrors: map[string]bool{
			"airly":              false,
			"search":             true,
			"searchInstallation": false,
			"installation":       false,
			"departures":         true,
		},
	}
	pb.RegisterKrkStopsServer(grpcServer, &server)
	grpc_prometheus.Register(grpcServer)
	log.Fatal(grpcServer.Serve(lis))
}
