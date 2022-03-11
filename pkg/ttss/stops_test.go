package ttss

import (
	"testing"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/matryer/is"
)

func TestStops(t *testing.T) {
	is := is.New(t)
	stopsC, errC := GetAllStops(testEndpoints)
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
