package ttss

import (
	"testing"
)

func TestGetAllStops(t *testing.T) {
	maps, err := GetAllStops()
	if err != nil {
		t.Fatal(err)
	}
	if len(maps) < 300 {
		t.Error("Got less than 300 stops!")
	}
}
