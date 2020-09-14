package main

import (
	"bufio"
	"log"
	"net/http"
	"os"

	"github.com/PiotrKozimor/krk-stops-backend-golang/krkstops"
	"github.com/spf13/cobra"
)

func init() {
	nonInteractive = *updateCmd.Flags().BoolP("interactive-skip", "i", false, "do not ask before updating stops")
}

var (
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update stops in Redis and suggestions in Redisearch",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)
			app := krkstops.App{}
			app.HTTPClient = &http.Client{}
			app.RedisAutocompleter = redisearchClient
			app.RedisClient = redisClient
			stops, err := app.GetAllStops()
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
