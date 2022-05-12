package main

import (
	"log"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/pkg/ttss"
	"github.com/spf13/cobra"
)

var (
	stopId  uint
	rootCmd = &cobra.Command{
		Use:  "ttssctl",
		Long: `ttssctl queries stops and departures from TTSS API of MPK in Cracow.`,
	}
	depsCmd = &cobra.Command{
		Use:   "deps",
		Short: "Query departures from given stop for bus and tram",
		Run: func(cmd *cobra.Command, args []string) {
			depsC, errC := ttss.GetDepartures(endpoints, stopId)
			pprint := pb.NewPrettyPrint(cmd)
			for dep := range depsC {
				pprint.Departures(getDepsPSlice(dep))
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
			stopsC, errC := ttss.GetAllStops(endpoints)
			cmd.Print()
			pprint := pb.NewPrettyPrint(cmd)
			for stop := range stopsC {
				pprint.Stops(getStopsPSlice(stop))
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
			deps, err := endpoints[pb.Endpoint_BUS].GetDepartures(stopId)
			if err != nil {
				log.Fatal(err)
			}
			pprint := pb.NewPrettyPrint(cmd)
			pprint.Departures(getDepsPSlice(deps))
		},
	}
	tramCmd = &cobra.Command{
		Use:   "tram",
		Short: "Query tram departures from stop",
		Run: func(cmd *cobra.Command, args []string) {
			deps, err := endpoints[pb.Endpoint_TRAM].GetDepartures(stopId)
			if err != nil {
				log.Fatal(err)
			}
			pprint := pb.NewPrettyPrint(cmd)
			pprint.Departures(getDepsPSlice(deps))
		},
	}
	endpoints = ttss.KrkStopsEndpoints
)

func getDepsPSlice(deps []pb.Departure) []*pb.Departure {
	depsP := make([]*pb.Departure, len(deps))
	for i := range deps {
		depsP[i] = &deps[i]
	}
	return depsP
}

func getStopsPSlice(stops []pb.Stop) []*pb.Stop {
	stopsP := make([]*pb.Stop, len(stops))
	for i := range stops {
		stopsP[i] = &stops[i]
	}
	return stopsP
}

func init() {
	depsCmd.PersistentFlags().UintVar(&stopId, "id", 610, "id of stop to query departures from")
	depsCmd.AddCommand(busCmd, tramCmd)
	rootCmd.AddCommand(stopsCmd, depsCmd)
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
