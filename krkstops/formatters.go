package krkstops

import (
	"fmt"
	"os"
	"text/tabwriter"

	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
)

var airlyWriter *tabwriter.Writer
var stopsWriter *tabwriter.Writer
var departuresWriter *tabwriter.Writer

func init() {
	airlyWriter = tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	stopsWriter = airlyWriter

}

// PrintAirly pretty with tabwriter
func PrintAirly(airly *pb.Airly) {
	fmt.Fprintf(airlyWriter, "CAQI\tHUMIDITY[%%]\tTEMP [°C]\t\n")
	fmt.Fprintf(airlyWriter, "%d\t%d\t%2.1f\t\n", airly.Caqi, airly.Humidity, airly.Temperature)
	airlyWriter.Flush()
}

// PrintStops pretty with tabwriter
func PrintStops(stops *StopsMap) {
	for shortName, name := range *stops {
		fmt.Fprintf(stopsWriter, "%s\t%s\t\n", shortName, name)
	}
	stopsWriter.Flush()
}

// PrintDepartures pretty with tabwriter
func PrintDepartures(w *tabwriter.Writer, deps []pb.Departure) {
	for _, dep := range deps {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t\n", dep.PatternText, dep.Direction, dep.RelativeTime, dep.PlannedTime)
	}
	w.Flush()
}
