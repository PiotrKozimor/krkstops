package cmd

import (
	gt "github.com/PiotrKozimor/krkstops/gtfs"
	"github.com/artonge/go-gtfs"
	"github.com/spf13/cobra"
)

var stopsCmd = &cobra.Command{
	Use:   "stops [gtfs folder]",
	Short: "Show inconsistencies in stops naming between GTFS and TTSS",
	Long:  `bus and tram folder is expected to exist in [gtfs folder]. Generate them with download_gtfs.sh script`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := gtfs.LoadSplitted(args[0], nil)
		handle(err)
		gt.CheckStops(data...)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(stopsCmd)
}
