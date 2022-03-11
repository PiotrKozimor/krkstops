package ttss

import (
	"context"
	"testing"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/test/mock"
	"github.com/matryer/is"
)

func TestStops(t *testing.T) {
	is := is.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go mock.Ttss(ctx)
	time.Sleep(10 * time.Millisecond)
	stopsC, errC := GetAllStops(KrkStopsEndpoints)
	for s := range stopsC {
		switch s[0].Type {
		case pb.Endpoint_BUS:
			is.Equal(2, len(s))
		case pb.Endpoint_TRAM:
			is.Equal(1, len(s))
		}
	}
	for err := range errC {
		is.NoErr(err)
	}
}
