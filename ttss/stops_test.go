package ttss

import (
	"context"
	"testing"
	"time"

	"github.com/PiotrKozimor/krkstops/mock"
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/stretchr/testify/assert"
)

func TestStops(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go mock.Ttss(ctx)
	time.Sleep(10 * time.Millisecond)
	stopsC, errC := GetAllStops(KrkStopsEndpoints)
	for s := range stopsC {
		switch s[0].Type {
		case pb.Endpoint_BUS:
			assert.Equal(t, 2, len(s))
		case pb.Endpoint_TRAM:
			assert.Equal(t, 1, len(s))
		}
	}
	for err := range errC {
		assert.NoError(t, err)
	}
}
