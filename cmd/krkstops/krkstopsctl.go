package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/PiotrKozimor/krk-stops-backend-golang/krkstops"
	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
	"github.com/spf13/cobra"
	grpc "google.golang.org/grpc"
)

var (
	rootCmd = &cobra.Command{
		Use:  "krkstopsctl",
		Long: `krkstopsctl interacts with krk-stops.pl backend`,
	}
	airlyCmd = &cobra.Command{
		Use:   "airly",
		Short: "Query air quality, temperature and humidity from airly installation",
		Run: func(cmd *cobra.Command, args []string) {
			initClient()
			airly, err := client.GetAirly(ctx, &pb.Installation{Id: airlyId})
			if err != nil {
				log.Fatalf("fail to get airly: %v", err)
			}
			krkstops.PrintAirly(airly)
		},
	}
	depsCmd = &cobra.Command{
		Use:   "deps",
		Short: "Query departures from given stop",
		Run: func(cmd *cobra.Command, args []string) {
			initClient()
			depsStrem, err := client.GetDepartures(ctx, &pb.Stop{ShortName: fmt.Sprint(stopId)})
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
		},
	}
	stopsCmd = &cobra.Command{
		Use:   "stops",
		Short: "Search stops",
		Run: func(cmd *cobra.Command, args []string) {
			initClient()
			stopsStream, err := client.SearchStops(ctx, &pb.StopSearch{Query: args[0]})
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
		},
		Args: cobra.ExactArgs(1),
	}

	airlyId  int32
	stopId   int32
	query    string
	conn     grpc.ClientConn
	client   pb.KrkStopsClient
	ctx      context.Context
	endpoint string
)

func initClient() {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client = pb.NewKrkStopsClient(conn)
	ctx, _ = context.WithTimeout(context.Background(), 1*time.Second)
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "krk-stops.pl:8080", "backend url address")
	depsCmd.Flags().Int32Var(&stopId, "id", 610, "stop id, find it by using stops command")
	airlyCmd.Flags().Int32Var(&airlyId, "id", 9914, "installation id")
	rootCmd.AddCommand(depsCmd)
	rootCmd.AddCommand(airlyCmd)
	rootCmd.AddCommand(stopsCmd)
}

func krkstopsUsage() {
	fmt.Printf(`Query KrkStops API
	deps           	fetch departures
	stops [QEURY]  	search stops by query
	airly   		air quality for given installation
`)
}

func main() {
	rootCmd.Execute()
}
