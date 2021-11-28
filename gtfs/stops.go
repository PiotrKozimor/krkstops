package gtfs

import (
	"log"

	"github.com/PiotrKozimor/krkstops/ttss"
	"github.com/artonge/go-gtfs"
)

func CheckStops(data ...*gtfs.GTFS) {
	ttssStops := make(map[string]bool, 1000)
	gtfsHeadsigns := make(map[string]bool, 100)
	stopsC, errC := ttss.GetAllStops(ttss.KrkStopsEndpoints)
	for stops := range stopsC {
		for i := range stops {
			ttssStops[stops[i].Name] = false
		}
	}
	for err := range errC {
		log.Panic(err)
	}
	for _, d := range data {
		for _, trip := range d.Trips {
			gtfsHeadsigns[trip.Headsign] = true
		}
	}
	for headsign := range gtfsHeadsigns {
		if _, ok := ttssStops[headsign]; !ok {
			log.Printf("not found in TTSS: %s", headsign)
		}
	}
	// for stop := range ttssStops {
	// 	if strings.Contains(stop, "Tow") {
	// 		log.Print(stop)
	// 	}
	// 	if !ttssStops[stop] {
	// 		log.Printf("not found in GTFS: %s", stop)

	// 	}
	// }
}
