package main

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)

}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete stops in Redis and suggestions in Redisearch",
	Run: func(cmd *cobra.Command, args []string) {
		initializeRedisClients()
		_, err := redisClient.FlushAll().Result()
		if err != nil {
			log.Fatal(err)
		}
		print("Stops deleted.\n")
	},
}
