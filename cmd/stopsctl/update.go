package main

import (
	"github.com/PiotrKozimor/krkstops/pkg/ttss"
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
			err = s.Update(ttss.BusEndpoint, ttss.TramEndpoint)
			handle(err)
			println("stops updated sucessfully.")
		},
	}
)
