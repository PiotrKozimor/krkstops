package ttss

import (
	"testing"

	"github.com/matryer/is"
)

func TestStops(t *testing.T) {
	is := is.New(t)
	cli := NewClient(Bus)
	stops, err := cli.GetAllStops()
	is.NoErr(err)
	t.Log(stops)
}

func TestTramStops(t *testing.T) {
	is := is.New(t)
	bus := NewClient(Bus)
	tram := NewClient(Tram)

	stops, err := tram.GetAllStops()
	is.NoErr(err)
	m := make(map[uint]string, len(stops))
	for _, stop := range stops {
		m[stop.Id] = stop.Name
	}

	stops, err = bus.GetAllStops()
	is.NoErr(err)
	for _, stop := range stops {
		m[stop.Id] = stop.Name
		delete(m, stop.Id)
	}

	t.Log(m)
}
