package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PiotrKozimor/krk-stops-backend-golang/krkstops"
	"github.com/spf13/cobra"
)

func init() {
	interactive = *updateCmd.Flags().BoolP("interactive", "i", true, "ask before updating stops")
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
			fmt.Println("New stops:")
			krkstops.PrintStops(&newStops)
			fmt.Println("Old stops:")
			krkstops.PrintStops(&oldStops)
			if interactive {
				fmt.Println("Apply new changes?[y/N] ")
				text, _ := reader.ReadString('\n')
				if text != "y\n" {
					os.Exit(0)
				}
			}
			err = app.UpdateSuggestionsAndRedis(newStops, oldStops)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	interactive bool
)
