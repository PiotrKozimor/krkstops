package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/PiotrKozimor/krkstops"
	"github.com/PiotrKozimor/krkstops/ttss"
	"github.com/spf13/cobra"
)

func init() {
	updateCmd.Flags().BoolVarP(&nonInteractive, "interactive-skip", "i", false, "do not ask before updating stops")
	rootCmd.AddCommand(updateCmd)
}

func ppStops(stops ttss.Stops) {
	for id, name := range stops {
		fmt.Fprintf(pp, "%d\t%s\t\n", id, name)
	}
	pp.Flush()
}

var (
	pp        = tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update stops in Redis and suggestions in Redisearch",
		Run: func(cmd *cobra.Command, args []string) {
			initializeRedisClients()
			reader := bufio.NewReader(os.Stdin)
			app := krkstops.Clients{
				Redis:              redisClient,
				RedisAutocompleter: redisearchClient,
			}
			stops := make(ttss.Stops)
			stopsC, errC := ttss.GetAllStops(ttss.KrkStopsEndpoints)
			for err := range errC {
				log.Fatal(err)
			}
			for stop := range stopsC {
				for _, s := range stop {
					stops[s.Id] = s.Name
				}
			}
			newStops, oldStops, err := app.CompareStops(stops)
			if err != nil {
				log.Fatal(err)
			}
			print("New stops:\n")
			ppStops(newStops)
			print("Old stops:\n")
			ppStops(oldStops)
			if !nonInteractive {
				print("Apply new changes?[y/N]\n")
				text, _ := reader.ReadString('\n')
				if text != "y\n" {
					os.Exit(0)
				}
			}
			err = app.UpdateSuggestionsAndRedis(newStops, oldStops)
			if err != nil {
				log.Fatal(err)
			}
			print("Stops updates sucessfully.\n")

		},
	}
	nonInteractive bool
)
