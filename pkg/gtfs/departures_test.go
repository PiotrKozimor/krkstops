package gtfs

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetDepartures(t *testing.T) {
	u, name, source, updates := Tram, "tram", "T", "testdata/TripUpdates_T.pb"
	stopName := "Ruczaj"

	d := NewDepartures(u, name)
	err := d.Init(func(file string) (*csv.Reader, error) {
		return mustRead(t, source, file), nil
	})
	require.NoError(t, err)

	b, err := os.ReadFile(updates)
	require.NoError(t, err)
	feed, err := ParseFeed(b)
	require.NoError(t, err)
	d.updates, err = d.ParseStopUpdates(feed)
	require.NoError(t, err)

	stopIds := d.Stops[stopName]

	for u, d := range d.updates {
		if slices.Contains(stopIds, u.stopId) {
			log.Printf("stopid %d\ttripid %d\tdelay %f", u.stopId, u.tripId, float64(d)/60)
		}
	}

	log := func(departures []Departure, headsigns []RouteHeadsign) {
		for _, dep := range departures {
			t.Logf("%v", dep)
		}
		for _, h := range headsigns {
			t.Logf("%v", h)
		}
	}
	log(d.Get("Ruczaj",
		time.Now().In(Location),
		120,
	))
}

func mustRead(t *testing.T, source, file string) *csv.Reader {
	b, err := os.ReadFile(fmt.Sprintf("testdata/GTFS_KRK_%s/%s", source, file))
	require.NoError(t, err)
	return csv.NewReader(bytes.NewBuffer(b))
}
