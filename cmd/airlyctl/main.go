package main

import (
	"log"

	"github.com/PiotrKozimor/krkstops/airly"
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:  "airlyctl",
		Long: `airlyctl queries CAQI, temperature and humidity from given Airly installation. Please provide API key via AIRLY_KEY environmental variable.`,
		Run: func(cmd *cobra.Command, args []string) {
			inst := pb.Installation{}
			inst.Id = id
			airly, err := airly.Api.GetAirly(&inst)
			if err != nil {
				log.Fatal(err)
			}
			pp := pb.NewPrettyPrint()
			pp.Airly(airly)
		},
	}
	id int32
)

func init() {
	rootCmd.Flags().Int32Var(&id, "id", 9895, "id of installation to query. Find it from map on https://airly.eu/map/pl/")
}

func main() {
	rootCmd.Execute()
}
