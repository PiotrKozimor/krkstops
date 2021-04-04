package krkstops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_scoreByTotalDepartures(t *testing.T) {
	assert.Equal(t, float64(1), scoreByTotalDepartures(0))
	assert.Equal(t, float64(5), scoreByTotalDepartures(64))
}
