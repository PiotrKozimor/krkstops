package main

import (
	"log"
	"net/http"

	"github.com/PiotrKozimor/krk-stops-backend-golang/krkstops"
	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:  "airlyctl",
		Long: `airlyctl queries CAQI, temperature and humidity from given Airly installation. Please provide API key via AIRLY_KEY environmental variable.`,
		Run: func(cmd *cobra.Command, args []string) {
			app := krkstops.App{}
			app.HTTPClient = &http.Client{}
			inst := pb.Installation{}
			inst.Id = id
			airly, err := app.GetAirly(&inst)
			if err != nil {
				log.Fatal(err)
			}
			krkstops.PrintAirly(&airly)
		},
	}
	id int32
)

func init() {
	rootCmd.Flags().Int32Var(&id, "id", 9914, "id of installation to query. Find it from map on https://airly.eu/map/pl/")
}

func main() {
	rootCmd.Execute()
}
