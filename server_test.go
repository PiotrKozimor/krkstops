package krkstops

import (
	"testing"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/matryer/is"
)

var (
	testInstallation = pb.Installation{
		Id: 8077,
	}
	testStop = pb.Stop{
		Id: 610,
	}
)

func TestGetAirly(t *testing.T) {
	is := is.New(t)
	c := mustLocalKrkStopsClient(is)
	airly, err := c.GetAirly(ctx, &testInstallation)
	is.NoErr(err)
	is.Equal(airly.Caqi, int32(12))
	is.Equal(airly.Humidity, int32(46))
	is.Equal(airly.Temperature, float32(12.33))
	is.Equal(airly.Color, uint32(7063846))
}

func TestGetDepartures2(t *testing.T) {
	is := is.New(t)
	c := mustLocalKrkStopsClient(is)
	deps, err := c.GetDepartures2(ctx, &testStop)
	is.NoErr(err)
	is.Equal(len(deps.Departures), 21)
	deps, err = c.GetDepartures2(ctx, &testStop)
	is.NoErr(err)
	is.Equal(len(deps.Departures), 21)
}

func TestSearchStops2(t *testing.T) {
	is := is.New(t)
	c := mustLocalKrkStopsClient(is)
	mustSearchEqual(is, "ma", "Rondo Matecznego", c)
	mustSearchEqual(is, "cza", "Czarnowiejska", c)
}

func mustSearchEqual(is *is.I, query, stopName string, cli pb.KrkStopsClient) {
	stops, err := cli.SearchStops2(ctx, &pb.StopSearch{
		Query: "ma",
	})
	is.NoErr(err)
	is.Equal(len(stops.Stops), 1)
	is.Equal(stops.Stops[0].Name, "Rondo Matecznego")
}
