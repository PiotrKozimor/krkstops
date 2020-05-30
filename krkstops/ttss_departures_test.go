package krkstops

import (
	"net/http"
	"testing"

	krk_stops "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
)

func TestTTSSDepartures(test *testing.T) {
	app := App{}
	app.HTTPClient = &http.Client{}
	departures, err := app.GetStopDepartures(&krk_stops.Stop{ShortName: "2688"})
	if err != nil {
		test.Error(err)
	}
	if len(departures) == 0 {
		test.Error("Got no departures!\n")
	}
}
