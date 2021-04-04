package ttss

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDepartures(t *testing.T) {
	departures, errC := GetDepartures(KrkStopsEndpoints, 2688)
	var cnt int
	for range departures {
		cnt++
	}
	for err := range errC {
		t.Fatal(err)
	}
	assert.Equal(t, cnt, 2)
}
