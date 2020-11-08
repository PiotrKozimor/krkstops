package main

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scoreCmd)
}

const STOPS_TO_SCORE = "stops.toScore"

var scoreCmd = &cobra.Command{
	Use:   "score",
	Short: "score stop suggestions in Redisearch (by number of departures from stop)",
	Run: func(cmd *cobra.Command, args []string) {
		// initializeRedisClients()
		// app := krkstops.App{}
		// app.RedisAutocompleter = redisearchClient
		// app.RedisClient = redisClient
		// exist, err := app.RedisClient.Exists(STOPS_TO_SCORE).Result()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// if exist == 1 {

		// }

	},
}
