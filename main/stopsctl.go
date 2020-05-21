package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/PiotrKozimor/krk-stops-backend-golang/krkstops"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v7"
)

func stopsUsage() {
	fmt.Printf(`Manipulate stops in Redis and stops suggestions in Redisearch:
	delete   delete all suggenstions
	update   fetch changes from TTSS API
	score    score stops
`)
}
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	scoreCmd := flag.NewFlagSet("score", flag.ExitOnError)
	flag.Usage = stopsUsage
	flag.Parse()
	w := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	reader := bufio.NewReader(os.Stdin)
	var interactive = updateCmd.Bool("i", false, "interactive mode")
	app := krkstops.App{}
	app.HTTPClient = &http.Client{}
	app.RedisClient = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	app.RedisAutocompleter = redisearch.NewAutocompleter("localhost:6379", "search-stops")
	switch os.Args[1] {
	case "delete":
		deleteCmd.Parse(os.Args[2:])
		_, err := app.RedisClient.FlushAll().Result()
		if err != nil {
			log.Fatal(err)
		}
	case "update":
		updateCmd.Parse(os.Args[2:])
		stops, err := app.GetAllStops()
		if err != nil {
			log.Fatal(err)
		}
		newStops, oldStops, err := app.CompareStops(stops)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("New stops:")
		krkstops.PrintStops(w, &newStops)
		fmt.Println("Old stops:")
		krkstops.PrintStops(w, &oldStops)
		if *interactive {
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
	case "score":
		scoreCmd.Parse(os.Args[2:])
		fmt.Print("Not implemented")
	}

}
