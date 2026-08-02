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

func (p *PrettyPrint) Stops(stops []string) {
	fmt.Fprintf(p, "NAME\n")
	for _, name := range stops {
		fmt.Fprintf(p, "%s\n", name)
	}
	p.Flush()
}

// PrettyPrint departures
func (p *PrettyPrint) Departures3(deps *pb.GetDeparturesResponse) {
	fmt.Fprintf(p, "DIR\tNAME\tPLAN\tUPDATE\n")
	for _, d := range deps.Departures {
		minutes := d.PlannedMinutesInDay % 60
		hours := (d.PlannedMinutesInDay - minutes) / 60
		if d.UpdatedSecondsInDay != 0 {
			fmt.Fprintf(p, "%d\t%d\t%02d:%02d\t%02d:%02d:%02d\n", d.DirectionId, d.RouteName, hours, minutes, d.UpdatedSecondsInDay/3600, d.UpdatedSecondsInDay/60%60, d.UpdatedSecondsInDay%60)
		} else {
			fmt.Fprintf(p, "%d\t%d\t%02d:%02d\n", d.DirectionId, d.RouteName, hours, minutes)

		}
	}
	fmt.Fprintf(p, "HEADSIGN\tROUTE NAME\tDIR ID\n")
	for _, h := range deps.Headsigns {
		fmt.Fprintf(p, "%s\t.\t.\n", h.Headsign)
		for _, r := range h.Routes {
			fmt.Fprintf(p, ".\t%d\t%d\n", r.RouteName, r.DirectionId)
		}
	}
	p.Flush()
}
