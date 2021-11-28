package cmd

import (
	"encoding/gob"
	"log"
	"os"

	gt "github.com/PiotrKozimor/krkstops/gtfs"
	"github.com/artonge/go-gtfs"
	"github.com/spf13/cobra"
)

func handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "gtfs [gtfs folder] [output file]",
	Short: "Parse GTFS files to departures grouping",
	Long:  `bus and tram folder is expected to exist in [gtfs folder]. Generate them with download_gtfs.sh script`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := gtfs.LoadSplitted(args[0], nil)
		handle(err)
		grouping, err := gt.ParseMany(data...)
		handle(err)
		f, err := os.OpenFile(args[1], os.O_RDWR|os.O_CREATE, 0755)
		handle(err)
		defer f.Close()
		enc := gob.NewEncoder(f)
		err = enc.Encode(grouping)
		handle(err)
	},
	Args: cobra.ExactArgs(2),
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
