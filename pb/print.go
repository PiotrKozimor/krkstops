package pb

import (
	"fmt"
	"text/tabwriter"

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
func (p PrettyPrint) Airly(airly *Airly) {
	fmt.Fprintf(p, "CAQI\tHUMIDITY[%%]\tTEMP [°C]\tCOLOR\t\n")
	fmt.Fprintf(p, "%d\t%d\t%2.1f\t%X\t\n", airly.Caqi, airly.Humidity, airly.Temperature, airly.Color)
	p.Flush()
}

// PrettyPrint stops
func (p *PrettyPrint) Stops(stops []*Stop) {
	fmt.Fprintf(p, "NO\tID\tNAME\n")
	for i := range stops {
		fmt.Fprintf(p, "%d\t%d\t%s\t\n", i, stops[i].Id, stops[i].Name)
	}
	p.Flush()
}

// PrettyPrint departures
func (p *PrettyPrint) Departures(deps []*Departure) {
	fmt.Fprintf(p, "ID\tNUMBER\tDIRECTION\tPLANNED\tRELATIVE\n")
	for _, dep := range deps {
		fmt.Fprintf(p, "%s\t%s\t%d\t%s\t%d\n", dep.PatternText, dep.Direction, dep.RelativeTime, dep.PlannedTime, dep.RelativeTime)
	}
	p.Flush()
}
