package ttss

import (
	"testing"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/matryer/is"
)

var testEndpoints = []Endpointer{
	Endpoint{
		Type: pb.Endpoint_BUS,
		URL:  "http://172.24.0.101:8071",
	},
	Endpoint{
		Type: pb.Endpoint_TRAM,
		URL:  "http://172.24.0.101:8070",
	},
}

func TestDepartures(t *testing.T) {
	is := is.New(t)
	departures, errC := GetDepartures(testEndpoints, 610)
	for d := range departures {
		switch d[0].Type {
		case pb.Endpoint_BUS:
			is.Equal(14, len(d))
			for i := range d {
				is.Equal(pb.Endpoint_BUS, d[i].Type)
			}
		case pb.Endpoint_TRAM:
			is.Equal(7, len(d))
			for i := range d {
				is.Equal(pb.Endpoint_TRAM, d[i].Type)
			}
		}
	}
	for err := range errC {
		is.NoErr(err)
	}
}
