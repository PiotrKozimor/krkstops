package ttss

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStops(t *testing.T) {
	stops, errC := GetAllStops(KrkStopsEndpoints)
	var cnt int
	for s := range stops {
		cnt++
		assert.Greater(t, len(s), 100)
		assert.Greater(t, len(s[0].Name), 2)
		assert.Greater(t, s[0].Id, uint32(0))
	}
	for err := range errC {
		t.Fatal(err)
	}
	assert.Equal(t, cnt, 2)
}
