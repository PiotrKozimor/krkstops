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

func airlyUsage() {
	fmt.Printf(`Query CAQI, temperature and humidity from airly installation. To find ID, please:
	1. visit https://airly.eu/map/pl/.
	2. Click on installation.
	3. ID will appear in URL.
Flags:
`)
	flag.PrintDefaults()
}

func main() {
	app := krkstops.App{}
	app.HTTPClient = &http.Client{}
	w := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var id = flag.Int("id", 9914, "ID of Airly installation to query")
	flag.Usage = airlyUsage
	flag.Parse()
	inst := pb.Installation{}
	inst.Id = int32(*id)
	airly, err := app.GetAirly(&inst)
	if err != nil {
		log.Fatal(err)
	}
	krkstops.PrintAirly(w, &airly)
}
