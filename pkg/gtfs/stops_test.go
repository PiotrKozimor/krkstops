package gtfs

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStopsIds(t *testing.T) {

	assert := func(stops Stops) {
		ids := make(map[uint32]bool)
		for _, s := range stops {
			for _, id := range s {
				assert.False(t, ids[id], "id: %s", id)
				ids[id] = true
			}
		}
	}

	stops, err := Bus.UnmarshalStops(mustReadStops(t, "A"))
	require.NoError(t, err)
	assert(stops)

	stops, err = Bus.UnmarshalStops(mustReadStops(t, "T"))
	require.NoError(t, err)
	assert(stops)
}

func mustReadStops(t *testing.T, source string) *csv.Reader {
	b, err := os.ReadFile(fmt.Sprintf("testdata/GTFS_KRK_%s/stops.txt", source))
	require.NoError(t, err)
	return csv.NewReader(bytes.NewBuffer(b))
}
