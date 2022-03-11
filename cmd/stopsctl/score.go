package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/PiotrKozimor/krkstops"
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	scoreCmd.Flags().StringVarP(&endpoint, "krkstops", "k", krkstops.ENDPOINT, "URL of krkstops endpoint")
	scoreCmd.Flags().DurationVarP(&sleep, "sleep", "s", time.Second, "sleep time between scoring stops")
	scoreCmd.Flags().BoolVarP(&restart, "restart", "r", false, "Discard existing scores")
	rootCmd.AddCommand(scoreCmd)

}

var (
	endpoint string
	sleep    time.Duration
	restart  bool
)

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
		ctx := cmd.Context()
		handle(err)
		if restart {
			err = score.RestartScoring(ctx)
			handle(err)
		}
		conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
		handle(err)
		client := pb.NewKrkStopsClient(conn)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		err = score.Score(ctx, c, client, sleep)
		handle(err)
		println("scoring finished")
	},
}
