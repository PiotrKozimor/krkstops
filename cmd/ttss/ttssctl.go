package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PiotrKozimor/krk-stops-backend-golang/krkstops"
	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
	"github.com/spf13/cobra"
)

var (
	app      = krkstops.App{HTTPClient: &http.Client{}}
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
				deps, err = app.GetStopDeparturesByURL(krkstops.BUS, &stop)
			} else if tramOnly {
				deps, err = app.GetStopDeparturesByURL(krkstops.TRAM, &stop)
			} else {
				deps, err = app.GetStopDepartures(&stop)
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
			var stops krkstops.StopsMap
			if busOnly {
				stops, err = app.GetAllStopsByURL(krkstops.BUS)
			} else if tramOnly {
				stops, err = app.GetAllStopsByURL(krkstops.TRAM)
			} else {
				stops, err = app.GetAllStops()
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
