package stops

import "testing"

func Test_scoreByTotalDepartures(t *testing.T) {
	if scoreByTotalDepartures(0) != 1 {
		t.Error()
	}
	if scoreByTotalDepartures(64) != 2 {
		t.Error()
	}
}
