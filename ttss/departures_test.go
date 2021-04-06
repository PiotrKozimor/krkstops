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
	departures, errC := GetDepartures([]Endpointer{
		Endpoint{
			URL:  "http://localhost:8080/internetservice",
			Type: pb.Endpoint_BUS},
		Endpoint{
			URL:  "http://localhost:8070/internetservice",
			Type: pb.Endpoint_TRAM},
	}, 610)
	for d := range departures {
		switch d[0].Type {
		case pb.Endpoint_BUS:
			assert.Equal(t, 14, len(d))
		case pb.Endpoint_TRAM:
			assert.Equal(t, 7, len(d))
		}
	}
	for err := range errC {
		assert.NoError(t, err)
	}
}
