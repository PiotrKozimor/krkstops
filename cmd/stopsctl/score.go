package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/PiotrKozimor/krkstops"
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	scoreCmd.Flags().StringVarP(&endpoint, "krkstops", "k", krkstops.ENDPOINT, "URL of krkstops endpoint")
	// scoreCmd.Flags().BoolVar(&restartScoring, "reset", false, "restart scoring")
	rootCmd.AddCommand(scoreCmd)

}

var endpoint string

// var restartScoring bool
var stop = false

const STOPS_TO_SCORE = "stops.toScore"

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var scoreCmd = &cobra.Command{
	Use: "score",
	Long: `score stop suggestions in Redisearch (by number of departures from stop). 
Stop scoring by sending interrupt signal (ctrl+C).`,
	Run: func(cmd *cobra.Command, args []string) {
		err := initializeDB()
		handle(err)
		conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
		handle(err)
		client := pb.NewKrkStopsClient(conn)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		err = db.Score(context.Background(), c, client)
		handle(err)
		println("scoring finished")
	},
}
