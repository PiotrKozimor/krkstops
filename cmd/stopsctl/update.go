package main

import (
	"bufio"
	"log"
	"os"

	"github.com/PiotrKozimor/krkstops"
	"github.com/PiotrKozimor/krkstops/stops"
	"github.com/PiotrKozimor/krkstops/ttss"
	"github.com/spf13/cobra"
)

func init() {
	updateCmd.Flags().BoolVarP(&nonInteractive, "interactive-skip", "i", false, "do not ask before updating stops")
	rootCmd.AddCommand(updateCmd)
}

var (
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update stops in Redis and suggestions in Redisearch",
		Run: func(cmd *cobra.Command, args []string) {
			initializeRedisClients()
			reader := bufio.NewReader(os.Stdin)
			app := stops.Clients{
				Redis:              redisClient,
				RedisAutocompleter: redisearchClient,
			}
			stops, err := ttss.GetAllStops()
			if err != nil {
				log.Fatal(err)
			}
			newStops, oldStops, err := app.CompareStops(stops)
			if err != nil {
				log.Fatal(err)
			}
			print("New stops:\n")
			krkstops.PrintStops(&newStops)
			print("Old stops:\n")
			krkstops.PrintStops(&oldStops)
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
