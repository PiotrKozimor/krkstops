package main

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var (
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update stops in Redis and suggestions in Redisearch",
		Run: func(cmd *cobra.Command, args []string) {
			err := initializeDB()
			handle(err)
			err = score.Update()
			handle(err)
			println("stops updated sucessfully.")
		},
	}
)
