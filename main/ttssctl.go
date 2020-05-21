package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/PiotrKozimor/krk-stops-backend-golang/krkstops"
	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
)

func ttssUsage() {
	fmt.Printf(`Fetch data from TTSS API:
	deps  fetch departures for given stop
	stops fetch all stops
`)
}

func main() {
	app := krkstops.App{}
	app.HTTPClient = &http.Client{}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	depsCmd := flag.NewFlagSet("deps", flag.ExitOnError)
	stopsCmd := flag.NewFlagSet("stops", flag.ExitOnError)
	var id = depsCmd.Int("id", 610, "shortName of stop to query")
	flag.Usage = ttssUsage
	flag.Parse()
	w := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	switch os.Args[1] {
	case "deps":
		depsCmd.Parse(os.Args[2:])
		stop := pb.Stop{}
		stop.ShortName = fmt.Sprint(*id)
		deps, err := app.GetStopDepartures(&stop)
		if err != nil {
			log.Fatal(err)
		}
		krkstops.PrintDepartures(w, deps)
	case "stops":
		stopsCmd.Parse(os.Args[2:])
		stops, err := app.GetAllStops()
		if err != nil {
			log.Fatal(err)
		}
		krkstops.PrintStops(w, &stops)
	}

}
