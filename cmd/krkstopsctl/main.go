package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/spf13/cobra"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/proto"
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
			pp.Stops3(stops.Stops)
		},
		Args: cobra.ExactArgs(1),
	}

	deps3Cmd = &cobra.Command{
		Use:   "deps3",
		Short: "Query departures from given stop",
		Run: func(cmd *cobra.Command, args []string) {
			b, _ := proto.Marshal(&pb.GetDepartures3Request{
				StopName: args[0],
			})
			resp, err := http.Post(endpoint+"/4", "", bytes.NewBuffer(b))
			handle(err)
			respB, err := io.ReadAll(resp.Body)
			handle(err)
			var r pb.GetDepartures3Response
			err = proto.Unmarshal(respB, &r)
			handle(err)
			pp := NewPrettyPrint(cmd)
			pp.Departures3(&r)
		},
	}
	stops3Cmd = &cobra.Command{
		Use:   "stops3",
		Short: "Search stops",
		Run: func(cmd *cobra.Command, args []string) {
			b, _ := proto.Marshal(&pb.SearchStops3Request{
				Query: args[0],
			})
			resp, err := http.Post(endpoint+"/3", "", bytes.NewBuffer(b))
			handle(err)
			log.Print(resp.StatusCode)
			respB, err := io.ReadAll(resp.Body)
			handle(err)
			var r pb.SearchStops3Response
			err = proto.Unmarshal(respB, &r)
			handle(err)
			pp := NewPrettyPrint(cmd)
			pp.Stops(r.Stops)
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
	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "http://localhost:8082", "backend url address")
	depsCmd.Flags().Int32Var(&stopId, "id", 610, "stop id, find it by using stops command")
	rootCmd.AddCommand(depsCmd)
	rootCmd.AddCommand(deps3Cmd)
	rootCmd.AddCommand(stopsCmd)
	rootCmd.AddCommand(stops3Cmd)

}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
