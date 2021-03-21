package krkstops

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/ttss"
)

type PrettyPrint struct {
	*tabwriter.Writer
}

func NewPrettyPrint() PrettyPrint {
	return PrettyPrint{
		tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0),
	}
}

// PrettyPrint Airly data
func (p PrettyPrint) Airly(airly *pb.Airly) {
	fmt.Fprintf(p, "CAQI\tHUMIDITY[%%]\tTEMP [°C]\tCOLOR\t\n")
	fmt.Fprintf(p, "%d\t%d\t%2.1f\t%X\t\n", airly.Caqi, airly.Humidity, airly.Temperature, airly.Color)
	p.Flush()
}

// PrettyPrint stops
func (p *PrettyPrint) Stops(stops []pb.Stop) {
	for i, stop := range stops {
		fmt.Fprintf(p, "%d\t%d\t%s\t\n", i, stop.Id, stop.Name)
	}
	p.Flush()
}

// PrettyPrint stopsMap
func (p *PrettyPrint) StopsMap(stops ttss.Stops) {
	for id, name := range stops {
		fmt.Fprintf(p, "%d\t%s\t\n", id, name)
	}
	p.Flush()
}

// PrettyPrint departures
func (p *PrettyPrint) Departures(deps []pb.Departure) {
	for _, dep := range deps {
		fmt.Fprintf(p, "%s\t%s\t%d\t%s\t\n", dep.PatternText, dep.Direction, dep.RelativeTime, dep.PlannedTime)
	}
	p.Flush()
}
