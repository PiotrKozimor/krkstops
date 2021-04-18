package airly

import (
	"context"
	"testing"

	"github.com/PiotrKozimor/krkstops/mock"
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/stretchr/testify/assert"
)

var mockEndp = Endpoint("http://localhost:8090")

func TestGetAirly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go mock.Airly(ctx)
	airly, err := mockEndp.GetAirly(&pb.Installation{Id: 8077})
	assert.NoError(t, err)
	assert.Equal(t, int32(12), airly.Caqi)
	assert.Equal(t, int32(46), airly.Humidity)
	assert.Equal(t, float32(12.33), airly.Temperature)
	assert.Equal(t, uint32(7063846), airly.Color)
	t.Logf("%+v", airly)
}

func TestFindAirlyInstallation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go mock.Airly(ctx)
	inst, err := mockEndp.NearestInstallation(&pb.InstallationLocation{Latitude: 50.0236288, Longitude: 19.942604799999998})
	assert.NoError(t, err)
	assert.Equal(t, int32(8077), inst.Id)
	t.Logf("%+v", inst)
}

func TestGetAirlyInstallation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go mock.Airly(ctx)
	inst, err := mockEndp.GetInstallation(8077)
	assert.NoError(t, err)
	t.Logf("%+v", inst)
	assert.Equal(t, float32(50.062008), inst.Latitude)
	assert.Equal(t, float32(19.940985), inst.Longitude)
}
