package main

import (
	"github.com/PiotrKozimor/krkstops/pkg/cache"
	"github.com/PiotrKozimor/krkstops/pkg/score"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "stopsctl",
		Short: "Manipulate krkstops stops data saved in Redis and stops suggestions saved in Redisearch.",
	}
	s   *score.Score
	uri string
)

func initializeDB() (err error) {
	s, err = score.NewScore(uri, cache.SUG)
	return
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&uri, "endpoint", "e", "localhost:6379", "URI of Redis endpoint")
}

func main() {
	rootCmd.Execute()
}
