package main

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"go.cozymore.dev/krkstops/pb"
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
			b, _ := proto.Marshal(&pb.GetDeparturesRequest{
				StopName: args[0],
			})
			resp, err := http.Post(endpoint+"/4", "", bytes.NewBuffer(b))
			handle(err)
			respB, err := io.ReadAll(resp.Body)
			handle(err)
			var r pb.GetDeparturesResponse
			err = proto.Unmarshal(respB, &r)
			handle(err)
			pp := NewPrettyPrint(cmd)
			pp.Departures3(&r)
		},
	}
	stopsCmd = &cobra.Command{
		Use:   "stops",
		Short: "Search stops",
		Run: func(cmd *cobra.Command, args []string) {
			b, _ := proto.Marshal(&pb.SearchStopsRequest{
				Query: args[0],
			})
			resp, err := http.Post(endpoint+"/3", "", bytes.NewBuffer(b))
			handle(err)
			log.Print(resp.StatusCode)
			respB, err := io.ReadAll(resp.Body)
			handle(err)
			var r pb.SearchStopsResponse
			err = proto.Unmarshal(respB, &r)
			handle(err)
			pp := NewPrettyPrint(cmd)
			pp.Stops(r.Stops)
		},
		Args: cobra.ExactArgs(1),
	}

	stopId   int32
	endpoint string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "http://localhost:8082", "backend url address")
	depsCmd.Flags().Int32Var(&stopId, "id", 610, "stop id, find it by using stops command")
	rootCmd.AddCommand(depsCmd)
	rootCmd.AddCommand(stopsCmd)

}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
