package krkstops

import (
	"fmt"
	"text/tabwriter"

	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
)

// PrintAirly pretty with tabwriter
func PrintAirly(w *tabwriter.Writer, airly *pb.Airly) {
	fmt.Fprintf(w, "CAQI\tHUMIDITY[%%]\tTEMP [°C]\t\n")
	fmt.Fprintf(w, "%d\t%d\t%2.1f\t\n", airly.Caqi, airly.Humidity, airly.Temperature)
	w.Flush()
}

// PrintStops pretty with tabwriter
func PrintStops(w *tabwriter.Writer, stops *StopsMap) {
	for shortName, name := range *stops {
		fmt.Fprintf(w, "%s\t%s\t\n", shortName, name)
	}
	w.Flush()
}

// PrintDepartures pretty with tabwriter
func PrintDepartures(w *tabwriter.Writer, deps []pb.Departure) {
	for _, dep := range deps {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t\n", dep.PatternText, dep.Direction, dep.RelativeTime, dep.PlannedTime)
	}
	w.Flush()
}
