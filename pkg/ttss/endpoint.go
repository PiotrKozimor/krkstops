package ttss

import (
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
	endpointsIds = [2]string{
		pb.Endpoint_BUS:  "bus",
		pb.Endpoint_TRAM: "tram",
	}
	KrkStopsEndpoints = []Endpointer{
		BusEndpoint,
		TramEndpoint,
	}
)
