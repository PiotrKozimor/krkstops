package gtfs

import (
	"log"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestGrouper(t *testing.T) {
	is := is.New(t)
	f, err := os.Open("krkstops.gob")
	is.NoErr(err)
	dut, err := NewGrouper(f)
	is.NoErr(err)
	dut.Dump(os.Stdout)
	// Rondo Grzegórzeckie
	mustBeGroup(is, dut, 365, map[string]string{
		"14": "Bronowice",
		"17": "Dworzec Towarowy",
		"19": "Dworzec Towarowy",
		"50": "Krowodrza Górka",
		"20": "Cichy Kącik",
		"9":  "Mistrzejowice",
		"49": "Tauron Arena Kraków Wieczysta",
		"12": "Tauron Arena Kraków Wieczysta",
	})

}

func mustBeGroup(i *is.I, dut *Grouper, stopId int, departures map[string]string) {
	for patterText, direction := range departures {
		d := Departure{
			StopId:      stopId,
			PatternText: patterText,
			Direction:   direction,
		}
		group := dut.Group(d)
		log.Printf("%+v\t%+v\n", group, d)
	}
}
