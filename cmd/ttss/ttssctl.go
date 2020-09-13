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
	app     = krkstops.App{HTTPClient: &http.Client{}}
	rootCmd = &cobra.Command{
		Use:  "ttssctl",
		Long: `ttss queries stops and departures from TTSS API of MPK in Cracow.`,
	}
	depsCmd = &cobra.Command{
		Use:   "deps",
		Short: "Query departures from given STOP_ID",
		Run: func(cmd *cobra.Command, args []string) {
			stop := pb.Stop{}
			stop.ShortName = fmt.Sprintf("%d", stopId)
			deps, err := app.GetStopDepartures(&stop)
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
			stops, err := app.GetAllStops()
			if err != nil {
				log.Fatal(err)
			}
			krkstops.PrintStops(&stops)
		},
	}
	stopId int32
)

func init() {
	depsCmd.Flags().Int32Var(&stopId, "id", 610, "id of stop to query departures from")
	rootCmd.AddCommand(stopsCmd)
	rootCmd.AddCommand(depsCmd)
}

func main() {
	// app := krkstops.App{}
	// app.HTTPClient = &http.Client{}
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	// depsCmd := flag.NewFlagSet("deps", flag.ExitOnError)
	// stopsCmd := flag.NewFlagSet("stops", flag.ExitOnError)
	// var id = depsCmd.Int("id", 610, "shortName of stop to query")
	// flag.Usage = ttssUsage
	// flag.Parse()
	// w := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	// switch os.Args[1] {
	// case "deps":
	// 	depsCmd.Parse(os.Args[2:])
	// 	stop := pb.Stop{}
	// 	stop.ShortName = fmt.Sprint(*id)
	// 	deps, err := app.GetStopDepartures(&stop)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	krkstops.PrintDepartures(deps)
	// case "stops":
	// 	stopsCmd.Parse(os.Args[2:])
	// 	stops, err := app.GetAllStops()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	krkstops.PrintStops(w, &stops)
	// }
	rootCmd.Execute()
}
