package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/PiotrKozimor/krkstops"
	"github.com/PiotrKozimor/krkstops/pb"
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
			defer cancel()
			airly, err := client.GetAirly(ctx, &pb.Installation{Id: airlyId})
			if err != nil {
				log.Fatalf("fail to get airly: %v", err)
			}
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
			depsStrem, err := client.GetDepartures(ctx, &pb.Stop{Id: uint32(stopId)})
			if err != nil {
				log.Fatalf("fail to get departures: %v", err)
			}
			out := cmd.OutOrStdout()
			for {
				departure, err := depsStrem.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintf(out, "%v\n", departure)
			}
		},
	}
	stopsCmd = &cobra.Command{
		Use:   "stops",
		Short: "Search stops",
		Run: func(cmd *cobra.Command, args []string) {
			initClient()
			defer cancel()
			stopsStream, err := client.SearchStops(ctx, &pb.StopSearch{Query: args[0]})
			if err != nil {
				log.Fatalf("fail to search stops: %v", err)
			}
			out := cmd.OutOrStdout()
			for {
				stop, err := stopsStream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintf(out, "%s %d\n", stop.Name, stop.Id)
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
