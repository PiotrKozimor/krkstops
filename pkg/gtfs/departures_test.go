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
	d := NewDepartures(Mobilis)
	err := d.Init(func(file string) (*csv.Reader, error) {
		return mustRead(t, "M", file), nil
	})
	require.NoError(t, err)

	b, err := os.ReadFile("testdata/TripUpdates_M.pb")
	require.NoError(t, err)
	feed, err := ParseFeed(b)
	require.NoError(t, err)
	d.updates, err = d.parseStopUpdates(feed)
	require.NoError(t, err)

	departures, headsigns := d.Get("Zachodnia",
		time.Date(2025, 12, 01, 20, 40, 0, 0, location),
		// time.Now().In(location),
		time.Date(2025, 12, 01, 21, 40, 0, 0, location),
		// time.Now().In(location).Add(time.Hour),
	)
	t.Log(departures)
	for _, dep := range departures {
		t.Logf("%+v", dep)
	}
	for _, h := range headsigns {
		t.Logf("%+v", h)
	}

}

func mustRead(t *testing.T, source, file string) *csv.Reader {
	b, err := os.ReadFile(fmt.Sprintf("testdata/GTFS_KRK_%s/%s", source, file))
	require.NoError(t, err)
	return csv.NewReader(bytes.NewBuffer(b))
}
