package ttss

import (
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
		URL:  os.Getenv("TTSS_BUS"),
		Type: pb.Endpoint_BUS,
	}
	TramEndpoint = Endpoint{
		URL:  os.Getenv("TTSS_TRAM"),
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
