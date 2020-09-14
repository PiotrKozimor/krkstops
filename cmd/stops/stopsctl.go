package main

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v7"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "stopsctl",
		Short: "Manipulate krkstops stops data saved in Redis and stops suggestions saved in Redisearch.",
	}
	redisClient      *redis.Client
	redisearchClient *redisearch.Autocompleter
)

func init() {
	url := rootCmd.Flags().StringP("endpoint", "e", "localhost:6379", "redis url address to connect to")
	redisClient = redis.NewClient(&redis.Options{Addr: *url})
	redisearchClient = redisearch.NewAutocompleter(*url, "search-stops")
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(scoreCmd)
	rootCmd.AddCommand(updateCmd)
}

func main() {
	rootCmd.Execute()
}
