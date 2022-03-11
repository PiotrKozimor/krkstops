package ttss

import (
	"context"
	"testing"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/test/mock"
	"github.com/matryer/is"
)

func TestDepartures(t *testing.T) {
	is := is.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go mock.Ttss(ctx)
	time.Sleep(10 * time.Millisecond)
	departures, errC := GetDepartures(KrkStopsEndpoints, 610)
	for d := range departures {
		switch d[0].Type {
		case pb.Endpoint_BUS:
			is.Equal(14, len(d))
			for _, dep := range d {
				is.Equal(pb.Endpoint_BUS, dep.Type)
			}
		case pb.Endpoint_TRAM:
			is.Equal(7, len(d))
			for _, dep := range d {
				is.Equal(pb.Endpoint_TRAM, dep.Type)
			}
		}
	}
	for err := range errC {
		is.NoErr(err)
	}
}
