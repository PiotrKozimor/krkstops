package main

import (
	"fmt"
	"text/tabwriter"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/spf13/cobra"
)

type PrettyPrint struct {
	*tabwriter.Writer
}

func NewPrettyPrint(cmd *cobra.Command) PrettyPrint {
	return PrettyPrint{
		tabwriter.NewWriter(cmd.OutOrStdout(), 1, 2, 2, ' ', 0),
	}
}

// PrettyPrint Airly data
func (p PrettyPrint) Measurement(m *pb.Measurement) {
	fmt.Fprintf(p, "CAQI\tHUMIDITY[%%]\tTEMP [°C]\tCOLOR\t\n")
	fmt.Fprintf(p, "%d\t%d\t%2.1f\t%X\t\n", m.Caqi, m.Humidity, m.Temperature, m.Color)
	p.Flush()
}

// PrettyPrint stops
func (p *PrettyPrint) Stops(stops []*pb.Stop) {
	fmt.Fprintf(p, "NO\tID\tNAME\n")
	for i := range stops {
		fmt.Fprintf(p, "%d\t%d\t%s\t\n", i, stops[i].Id, stops[i].Name)
	}
	p.Flush()
}

// PrettyPrint departures
func (p *PrettyPrint) Departures(deps []*pb.Departure) {
	fmt.Fprintf(p, "NO\tID\tDIRECTION\tPLANNED\tRELATIVE\n")
	for i, dep := range deps {
		fmt.Fprintf(p, "%d\t%s\t%s\t%s\t%d\n", i, dep.PatternText, dep.Direction, dep.PlannedTime, dep.RelativeTime)
	}
	p.Flush()
}
