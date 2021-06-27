package krkstops

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"google.golang.org/grpc"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func TestGrpc(t *testing.T) {
	conn, err := grpc.Dial(
		os.Getenv("PC_HOST"),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client := pb.NewKrkStopsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	a, err := client.GetAirly(ctx, &pb.Installation{
		Id: 9895,
	})
	log.Printf("a: %#+v\n", a)
	handle(err)
	scl, err := client.SearchStops(ctx, &pb.StopSearch{
		Query: "nor",
	})
	handle(err)
	for {
		stop, err := scl.Recv()
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("stop: %#+v\n", stop)
	}
	dcl, err := client.GetDepartures(ctx, &pb.Stop{
		Id: 2688,
	})
	handle(err)
	for {
		dep, err := dcl.Recv()
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("dep: %#+v\n", dep)
	}

}

func TestGrpcNoStream(t *testing.T) {
	conn, err := grpc.Dial(
		os.Getenv("PC_HOST"),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client := pb.NewKrkStopsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	a, err := client.GetAirly(ctx, &pb.Installation{
		Id: 9895,
	})
	log.Printf("a: %#+v\n", a)
	handle(err)
	stops, err := client.SearchStops2(ctx, &pb.StopSearch{
		Query: "nor",
	})
	handle(err)
	log.Printf("stops: %#+v\n", stops)
	deps, err := client.GetDepartures(ctx, &pb.Stop{
		Id: 2688,
	})
	handle(err)
	log.Printf("deps: %#+v\n", deps)
}
