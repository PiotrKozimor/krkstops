package ttss

import (
	"context"
	"testing"
	"time"

	"github.com/PiotrKozimor/krkstops/mock"
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/stretchr/testify/assert"
)

func TestDepartures(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go mock.Ttss(ctx)
	time.Sleep(10 * time.Millisecond)
	departures, errC := GetDepartures(KrkStopsEndpoints, 610)
	for d := range departures {
		switch d[0].Type {
		case pb.Endpoint_BUS:
			assert.Equal(t, 14, len(d))
			for _, dep := range d {
				assert.Equal(t, pb.Endpoint_BUS, dep.Type)
			}
		case pb.Endpoint_TRAM:
			assert.Equal(t, 7, len(d))
			for _, dep := range d {
				assert.Equal(t, pb.Endpoint_TRAM, dep.Type)
			}
		}
	}
	for err := range errC {
		assert.NoError(t, err)
	}
}
