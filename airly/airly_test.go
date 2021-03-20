package airly

import (
	"log"
	"testing"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/stretchr/testify/assert"
)

func TestGetAirly(t *testing.T) {
	airly, err := GetAirly(&pb.Installation{Id: 8077})
	assert.NoError(t, err)
	assert.Greater(t, airly.Humidity, int32(0))
	assert.Greater(t, airly.Caqi, int32(0))
	assert.Greater(t, airly.Color, uint32(0))
}

func TestFindAirlyInstallation(t *testing.T) {
	inst, err := NearestInstallation(&pb.InstallationLocation{Latitude: 50.0236288, Longitude: 19.942604799999998})
	assert.NoError(t, err)
	log.Print(inst)
}

func TestGetAirlyInstallation(t *testing.T) {
	inst, err := GetInstallation(8077)
	assert.NoError(t, err)
	log.Print(inst)
}
