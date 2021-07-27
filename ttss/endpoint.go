package ttss

import (
	"log"
	"os"

	"github.com/PiotrKozimor/krkstops/pb"
)

type Endpoint struct {
	URL  string
	Type pb.Endpoint
}

type Endpointer interface {
	GetDepartures(uint) ([]pb.Departure, error)
	GetAllStops() ([]pb.Stop, error)
	Id() string
}

var (
	BusEndpoint = Endpoint{
		URL:  "http://91.223.13.70",
		Type: pb.Endpoint_BUS,
	}
	TramEndpoint = Endpoint{
		URL:  "http://185.70.182.51",
		Type: pb.Endpoint_TRAM,
	}
	endpointsIds = map[Endpoint]string{
		BusEndpoint:  "bus",
		TramEndpoint: "tram",
	}
	KrkStopsEndpoints = []Endpointer{
		BusEndpoint,
		TramEndpoint,
	}
)

func init() {
	bus := os.Getenv("OVERRIDE_TTSS_BUS")
	if bus != "" {
		log.Printf("OVERRIDE_TTSS_BUS: %s\n", bus)
		BusEndpoint = Endpoint{
			URL:  bus,
			Type: pb.Endpoint_BUS,
		}
	}
	tram := os.Getenv("OVERRIDE_TTSS_TRAM")
	if tram != "" {
		log.Printf("OVERRIDE_TTSS_TRAM: %s\n", tram)
		TramEndpoint = Endpoint{
			URL:  tram,
			Type: pb.Endpoint_TRAM,
		}
	}
	KrkStopsEndpoints = []Endpointer{
		BusEndpoint,
		TramEndpoint,
	}
}
