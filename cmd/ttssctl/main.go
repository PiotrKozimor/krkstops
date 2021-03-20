package main

import (
	"log"

	"github.com/PiotrKozimor/krkstops"
	"github.com/PiotrKozimor/krkstops/ttss"
	"github.com/spf13/cobra"
)

var (
	busOnly  bool
	tramOnly bool
	rootCmd  = &cobra.Command{
		Use:  "ttssctl",
		Long: `ttssctl queries stops and departures from TTSS API of MPK in Cracow.`,
	}
	depsCmd = &cobra.Command{
		Use:   "deps",
		Short: "Query departures from given stop for bus and tram",
		Run: func(cmd *cobra.Command, args []string) {
			depsC, errC := ttss.GetDepartures(ttss.KrkStopsEndpoints, stopId)
			pprint := krkstops.NewPrettyPrint()
			for dep := range depsC {
				pprint.Departures(dep)
			}
			for err := range errC {
				log.Fatal(err)
			}
		},
	}
	stopsCmd = &cobra.Command{
		Use:   "stops",
		Short: "Query all stops in Cracov",
		Run: func(cmd *cobra.Command, args []string) {
			stopsC, errC := ttss.GetAllStops(ttss.KrkStopsEndpoints)
			pprint := krkstops.NewPrettyPrint()
			for stop := range stopsC {
				pprint.Stops(stop)
			}
			for err := range errC {
				log.Fatal(err)
			}
		},
	}
	busCmd = &cobra.Command{
		Use:   "bus",
		Short: "Query bus departures from stop",
		Run: func(cmd *cobra.Command, args []string) {
			deps, err := ttss.BusEndpoint.GetDepartures(stopId)
			if err != nil {
				log.Fatal(err)
			}
			pprint := krkstops.NewPrettyPrint()
			pprint.Departures(deps)
		},
	}
	tramCmd = &cobra.Command{
		Use:   "tram",
		Short: "Query tram departures from stop",
		Run: func(cmd *cobra.Command, args []string) {
			deps, err := ttss.TramEndpoint.GetDepartures(stopId)
			if err != nil {
				log.Fatal(err)
			}
			pprint := krkstops.NewPrettyPrint()
			pprint.Departures(deps)
		},
	}
	stopId uint
)

func init() {
	depsCmd.Flags().UintVar(&stopId, "id", 610, "id of stop to query departures from")
	depsCmd.AddCommand(busCmd, tramCmd)
	rootCmd.AddCommand(stopsCmd, depsCmd)
}

func main() {
	rootCmd.Execute()
}
