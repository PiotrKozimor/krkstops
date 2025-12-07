package gtfs

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetDepartures(t *testing.T) {
	d := NewDepartures(Bus)
	err := d.Init(func(file string) (*csv.Reader, error) {
		return mustRead(t, "A", file), nil
	})
	require.NoError(t, err)

	b, err := os.ReadFile("testdata/TripUpdates_A.pb")
	require.NoError(t, err)
	feed, err := ParseFeed(b)
	require.NoError(t, err)
	d.updates, err = d.ParseStopUpdates(feed)
	require.NoError(t, err)

	log := func(departures []Departure, headsigns []RouteHeadsign) {
		for _, dep := range departures {
			t.Logf("%+v", dep)
		}
		for _, h := range headsigns {
			t.Logf("%+v", h)
		}
	}
	log(d.Get("Zachodnia",
		time.Date(2025, 12, 8, 1, 1, 0, 0, Location),
		TwoHours,
	))
	log(d.Get("Teatr Słowackiego",
		time.Date(2025, 12, 8, 0, 0, 0, 0, Location),
		TwoHours,
	))
	log(d.Get("Teatr Słowackiego",
		time.Date(2025, 12, 7, 23, 50, 0, 0, Location),
		TwoHours,
	))

}

func mustRead(t *testing.T, source, file string) *csv.Reader {
	b, err := os.ReadFile(fmt.Sprintf("testdata/GTFS_KRK_%s/%s", source, file))
	require.NoError(t, err)
	return csv.NewReader(bytes.NewBuffer(b))
}
