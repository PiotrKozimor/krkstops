package main

import (
	"fmt"
	"log"

	"github.com/PiotrKozimor/krkstops"
	"github.com/PiotrKozimor/krkstops/pb"
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
		Short: "Query departures from given stop",
		Run: func(cmd *cobra.Command, args []string) {
			stop := pb.Stop{}
			stop.ShortName = fmt.Sprintf("%d", stopId)
			var err error
			var deps []pb.Departure
			if busOnly {
				deps, err = ttss.GetStopDeparturesByURL(ttss.BUS, &stop)
			} else if tramOnly {
				deps, err = ttss.GetStopDeparturesByURL(ttss.TRAM, &stop)
			} else {
				deps, err = ttss.GetStopDepartures(&stop)
			}
			if err != nil {
				log.Fatal(err)
			}
			krkstops.PrintDepartures(deps)
		},
	}
	stopsCmd = &cobra.Command{
		Use:   "stops",
		Short: "Query all stops in Cracov",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var stops ttss.Stops
			if busOnly {
				stops, err = ttss.GetAllStopsByURL(ttss.BUS)
			} else if tramOnly {
				stops, err = ttss.GetAllStopsByURL(ttss.TRAM)
			} else {
				stops, err = ttss.GetAllStops()
			}
			if err != nil {
				log.Fatal(err)
			}
			krkstops.PrintStops(&stops)
		},
	}
	stopId int32
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&busOnly, "bus", "b", false, "fetch bus TTSS API only")
	rootCmd.PersistentFlags().BoolVarP(&tramOnly, "tram", "t", false, "fetch tram TTSS API only")
	depsCmd.Flags().Int32Var(&stopId, "id", 610, "id of stop to query departures from")
	rootCmd.AddCommand(stopsCmd)
	rootCmd.AddCommand(depsCmd)
}

func main() {
	rootCmd.Execute()
}
