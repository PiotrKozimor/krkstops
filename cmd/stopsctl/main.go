package main

import (
	"context"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "stopsctl",
		Short: "Manipulate krkstops stops data saved in Redis and stops suggestions saved in Redisearch.",
	}
	redisClient      *redis.Client
	redisearchClient *redisearch.Autocompleter
	url              string
	ctx              = context.Background()
)

func initializeRedisClients() {
	redisClient = redis.NewClient(&redis.Options{Addr: url})
	redisearchClient = redisearch.NewAutocompleter(url, "search-stops")
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&url, "endpoint", "e", "localhost:6379", "redis url address to connect to")
}

func main() {
	rootCmd.Execute()
}
