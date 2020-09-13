package main

import (
	"log"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete stops in Redis and suggestions in Redisearch",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := redisClient.FlushAll().Result()
		if err != nil {
			log.Fatal(err)
		}
	},
}
