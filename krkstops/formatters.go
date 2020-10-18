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
	departuresWriter = airlyWriter
}

// PrintAirly pretty with tabwriter
func PrintAirly(airly *pb.Airly) {
	fmt.Fprintf(airlyWriter, "CAQI\tHUMIDITY[%%]\tTEMP [°C]\tCOLOR\t\n")
	fmt.Fprintf(airlyWriter, "%d\t%d\t%2.1f\t%s\t\n", airly.Caqi, airly.Humidity, airly.Temperature, airly.Color)
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
func PrintDepartures(deps []pb.Departure) {
	for _, dep := range deps {
		fmt.Fprintf(departuresWriter, "%s\t%s\t%d\t%s\t\n", dep.PatternText, dep.Direction, dep.RelativeTime, dep.PlannedTime)
	}
	departuresWriter.Flush()
}
