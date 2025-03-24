package main

import (
	"context"
	"log"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/spf13/cobra"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	depsCmd = &cobra.Command{
		Use:   "deps",
		Short: "Query departures from given stop",
		Run: func(cmd *cobra.Command, args []string) {
			initClient()
			defer cancel()
			deps, err := client.GetDepartures2(ctx, &pb.GetDepartures2Request{Id: uint32(stopId)})
			handle(err)
			pp := NewPrettyPrint(cmd)
			pp.Departures(deps.Departures)
		},
	}
	stopsCmd = &cobra.Command{
		Use:   "stops",
		Short: "Search stops",
		Run: func(cmd *cobra.Command, args []string) {
			initClient()
			defer cancel()
			stops, err := client.SearchStops2(ctx, &pb.SearchStops2Request{Query: args[0]})
			handle(err)
			pp := NewPrettyPrint(cmd)
			pp.Stops(stops.Stops)
		},
		Args: cobra.ExactArgs(1),
	}

	stopId   int32
	client   pb.KrkStopsClient
	ctx      context.Context
	endpoint string
	cancel   context.CancelFunc
)

func initClient() {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client = pb.NewKrkStopsClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "apikrkstops.cozymore.dev", "backend url address")
	depsCmd.Flags().Int32Var(&stopId, "id", 610, "stop id, find it by using stops command")
	rootCmd.AddCommand(depsCmd)
	rootCmd.AddCommand(stopsCmd)

}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
