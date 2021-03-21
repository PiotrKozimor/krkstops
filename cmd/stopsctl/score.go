package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/stops"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	scoreCmd.Flags().StringVar(&endpoint, "krk-endpoint", "krkstops.germanywestcentral.cloudapp.azure.com:8080", "url address of krk-stops backend to connect to")
	scoreCmd.Flags().BoolVar(&restartScoring, "reset", false, "restart scoring")
	rootCmd.AddCommand(scoreCmd)

}

var endpoint string
var restartScoring bool
var stop = false

const STOPS_TO_SCORE = "stops.toScore"

var scoreCmd = &cobra.Command{
	Use:   "score",
	Short: "score stop suggestions in Redisearch (by number of departures from stop)",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			<-c
			stop = true
		}()
		conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		client := pb.NewKrkStopsClient(conn)
		initializeRedisClients()
		cl := stops.Clients{
			Redis:              redisClient,
			RedisAutocompleter: redisearchClient,
		}
		if restartScoring {
			err = cl.RestartScoring()
		} else {
			err = cl.InitializeScoring()
		}
		if err == stops.ScoringInitialized {
			log.Println("Scoring is already initialized, continuing")
		} else if err != nil {
			log.Fatal(err)
		}
		for {
			if stop {
				log.Print("Exiting gracefully")
				break
			}
			stop, err := cl.GetStopToScore()
			if err != nil {
				log.Fatal(err)
			}
			if err == redis.Nil {
				log.Println("All stops scored")
				break
			}
			score, err := cl.ScoreStop(client, stop)
			if err != nil {
				log.Fatal(err)
			}
			cl.ApplyScore(score, stop)
			log.Printf("Score %.3f \t for stop: %s applied\n", score, stop)
			time.Sleep(time.Second)
		}
	},
}
