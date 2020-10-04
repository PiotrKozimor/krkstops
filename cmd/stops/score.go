package main

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scoreCmd)
}

var scoreCmd = &cobra.Command{
	Use:   "score",
	Short: "score stop suggestions in Redisearch (by number of departures from stop)",
	Run: func(cmd *cobra.Command, args []string) {
		print("Not implemented.\n")
	},
}
