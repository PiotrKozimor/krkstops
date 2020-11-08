package ttss

import (
	"testing"

	"github.com/PiotrKozimor/krkstops/pb"
)

func TestTTSSDepartures(test *testing.T) {
	departures, err := GetStopDepartures(&pb.Stop{ShortName: "2688"})
	if err != nil {
		test.Error(err)
	}
	if len(departures) == 0 {
		test.Error("Got no departures!\n")
	}
}
