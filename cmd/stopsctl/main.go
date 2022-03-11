package main

import (
	"github.com/PiotrKozimor/krkstops"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "stopsctl",
		Short: "Manipulate krkstops stops data saved in Redis and stops suggestions saved in Redisearch.",
	}
	score *krkstops.Score
	uri   string
)

func initializeDB() (err error) {
	score, err = krkstops.NewScore(uri, krkstops.SUG)
	return
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&uri, "endpoint", "e", "localhost:6379", "URI of Redis endpoint")
}

func main() {
	rootCmd.Execute()
}
