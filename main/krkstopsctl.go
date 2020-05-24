package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/PiotrKozimor/krk-stops-backend-golang/krkstops"
	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
	grpc "google.golang.org/grpc"
)

var conn grpc.ClientConn
var client pb.KrkStopsClient

func krkstopsUsage() {
	fmt.Printf(`Query KrkStops API
	deps           	fetch departures
	stops [QEURY]  	search stops by query
	airly   		air quality for given installation
`)
}

func main() {
	w := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	var endpoint string = "krk-stops.pl"
	depsCmd := flag.NewFlagSet("deps", flag.ExitOnError)
	// stopsCmd := flag.NewFlagSet("stops", flag.ExitOnError)
	airlyCmd := flag.NewFlagSet("airly", flag.ExitOnError)
	var airlyID = airlyCmd.Int("id", 9914, "ID of Airly installation to query")
	var depShortName = depsCmd.String("id", "610", "shortName of stop to query")
	flag.Usage = krkstopsUsage
	flag.Parse()
	endpointEnv := os.Getenv("ENDPOINT")
	if endpointEnv != "" {
		endpoint = endpointEnv
	}
	fmt.Println(endpoint)
	conn, err := grpc.Dial(fmt.Sprintf("%s:8080", endpoint), grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client := pb.NewKrkStopsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	switch os.Args[1] {
	case "airly":
		airlyCmd.Parse(os.Args[2:])
		airly, err := client.GetAirly(ctx, &pb.Installation{Id: int32(*airlyID)})
		if err != nil {
			log.Fatalf("fail to get airly: %v", err)
		}
		krkstops.PrintAirly(w, airly)
	case "deps":
		depsCmd.Parse(os.Args[2:])
		depsStrem, err := client.GetDepartures(ctx, &pb.Stop{ShortName: *depShortName})
		if err != nil {
			log.Fatalf("fail to get departures: %v", err)
		}
		for {
			departure, err := depsStrem.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			log.Println(departure)
		}
	case "stops":
		stopsStream, err := client.SearchStops(ctx, &pb.StopSearch{Query: os.Args[2]})
		if err != nil {
			log.Fatalf("fail to search stops: %v", err)
		}
		for {
			stop, err := stopsStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			log.Println(stop)
		}
	default:
		fmt.Print("Please provide correct command!")
	}
}
