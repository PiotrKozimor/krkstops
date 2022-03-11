package airly

import (
	"testing"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/matryer/is"
)

var testEndpoint = Endpoint("http://172.24.0.101:8072")

func TestGetAirly(t *testing.T) {
	is := is.New(t)
	airly, err := testEndpoint.GetAirly(&pb.Installation{Id: 8077})
	is.NoErr(err)
	is.Equal(int32(12), airly.Caqi)
	is.Equal(int32(46), airly.Humidity)
	is.Equal(float32(12.33), airly.Temperature)
	is.Equal(uint32(7063846), airly.Color)
}

func TestFindAirlyInstallation(t *testing.T) {
	is := is.New(t)
	inst, err := testEndpoint.NearestInstallation(&pb.InstallationLocation{Latitude: 50.0236288, Longitude: 19.942604799999998})
	is.NoErr(err)
	is.Equal(int32(8077), inst.Id)
}

func TestGetAirlyInstallation(t *testing.T) {
	is := is.New(t)
	inst, err := testEndpoint.GetInstallation(8077)
	is.NoErr(err)
	is.Equal(float32(50.062008), inst.Latitude)
	is.Equal(float32(19.940985), inst.Longitude)
}
