package ttss

import "github.com/PiotrKozimor/krkstops/pb"

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
		URL:  "http://91.223.13.70/internetservice",
		Type: pb.Endpoint_BUS,
	}
	TramEndpoint = Endpoint{
		URL:  "http://185.70.182.51/internetservice",
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
