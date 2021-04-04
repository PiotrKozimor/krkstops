package main

import (
	"bufio"
	"log"
	"os"

	"github.com/PiotrKozimor/krkstops"
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
			pp := krkstops.NewPrettyPrint()
			print("New stops:\n")
			pp.StopsMap(newStops)
			print("Old stops:\n")
			pp.StopsMap(oldStops)
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
