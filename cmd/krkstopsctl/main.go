package main

import (
	"context"
	"log"
	"time"

	"github.com/PiotrKozimor/krkstops"
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/spf13/cobra"
	grpc "google.golang.org/grpc"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

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
			defer cancel()
			airly, err := client.GetAirly(ctx, &pb.Installation{Id: airlyId})
			handle(err)
			pp := pb.NewPrettyPrint(cmd)
			pp.Airly(airly)
		},
	}
	depsCmd = &cobra.Command{
		Use:   "deps",
		Short: "Query departures from given stop",
		Run: func(cmd *cobra.Command, args []string) {
			initClient()
			defer cancel()
			deps, err := client.GetDepartures2(ctx, &pb.Stop{Id: uint32(stopId)})
			handle(err)
			pp := pb.NewPrettyPrint(cmd)
			pp.Departures(deps.Departures)
		},
	}
	stopsCmd = &cobra.Command{
		Use:   "stops",
		Short: "Search stops",
		Run: func(cmd *cobra.Command, args []string) {
			initClient()
			defer cancel()
			stops, err := client.SearchStops2(ctx, &pb.StopSearch{Query: args[0]})
			handle(err)
			pp := pb.NewPrettyPrint(cmd)
			pp.Stops(stops.Stops)
		},
		Args: cobra.ExactArgs(1),
	}

	airlyId  int32
	stopId   int32
	client   pb.KrkStopsClient
	ctx      context.Context
	endpoint string
	cancel   context.CancelFunc
)

func initClient() {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client = pb.NewKrkStopsClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", krkstops.ENDPOINT, "backend url address")
	depsCmd.Flags().Int32Var(&stopId, "id", 610, "stop id, find it by using stops command")
	airlyCmd.Flags().Int32Var(&airlyId, "id", 9895, "installation id")
	rootCmd.AddCommand(depsCmd)
	rootCmd.AddCommand(airlyCmd)
	rootCmd.AddCommand(stopsCmd)

}

func main() {
	rootCmd.Execute()
}
