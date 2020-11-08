package airly

import (
	"log"
	"strings"
	"testing"

	"github.com/PiotrKozimor/krkstops/pb"
)

func TestGetAirly(test *testing.T) {
	airly, err := GetAirly(&pb.Installation{Id: 8077})
	if err != nil {
		test.Error(err)
	}
	if airly.Humidity < 0 {
		test.Errorf("Humdidity %d < 0", airly.Humidity)
	}
	if airly.Caqi < 0 {
		test.Errorf("CAQI %d < 0", airly.Caqi)
	}
	if strings.Contains(airly.Color, "#") != true {
		test.Errorf("Got improper color string: %s", airly.Color)
	}
}

func TestFindAirlyInstallation(test *testing.T) {
	inst, err := FindAirlyInstallation(&pb.InstallationLocation{Latitude: 50.0236288, Longitude: 19.942604799999998})
	if err != nil {
		test.Error(err)
	}
	log.Print(inst)
}

func TestGetAirlyInstallation(test *testing.T) {
	inst, err := GetAirlyInstallation(&pb.Installation{
		Id: 8077,
	})
	if err != nil {
		test.Error(err)
	}
	log.Print(inst)
}
