package airly

import (
	"context"
	"testing"

	"github.com/PiotrKozimor/krkstops/mock"
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/matryer/is"
)

func TestGetAirly(t *testing.T) {
	is := is.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go mock.Airly(ctx)
	airly, err := DefaultEndpoint.GetAirly(&pb.Installation{Id: 8077})
	is.NoErr(err)
	is.Equal(int32(12), airly.Caqi)
	is.Equal(int32(46), airly.Humidity)
	is.Equal(float32(12.33), airly.Temperature)
	is.Equal(uint32(7063846), airly.Color)
}

func TestFindAirlyInstallation(t *testing.T) {
	is := is.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go mock.Airly(ctx)
	inst, err := DefaultEndpoint.NearestInstallation(&pb.InstallationLocation{Latitude: 50.0236288, Longitude: 19.942604799999998})
	is.NoErr(err)
	is.Equal(int32(8077), inst.Id)
}

func TestGetAirlyInstallation(t *testing.T) {
	is := is.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go mock.Airly(ctx)
	inst, err := DefaultEndpoint.GetInstallation(8077)
	is.NoErr(err)
	is.Equal(float32(50.062008), inst.Latitude)
	is.Equal(float32(19.940985), inst.Longitude)
}
