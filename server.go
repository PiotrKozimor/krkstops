package krkstops

import (
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/pkg/airly"
	"github.com/PiotrKozimor/krkstops/pkg/search"
	"github.com/PiotrKozimor/krkstops/pkg/ttss"
)

type KrkStopsServer struct {
	pb.UnimplementedKrkStopsServer
	airlyCache  cache[airly.Measurement]
	depsCache   cache[[]typedDeparture]
	airlyCli    airlyClient
	ttssBusCli  ttssClient
	ttssTramCli ttssClient
	searchCli   *search.Search
}

type airlyClient interface {
	GetMeasurement(installationId uint) (airly.Measurement, error)
	GetInstallation(id uint) (airly.Installation, error)
	NearestInstallation(lat, lon float32) (airly.Installation, error)
}

type ttssClient interface {
	GetDepartures(uint) ([]ttss.Departure, error)
	// GetAllStops() ([]Stop, error)
}

func NewServer() (*KrkStopsServer, error) {
	s, err := search.New()
	if err != nil {
		return nil, err
	}
	return &KrkStopsServer{
		airlyCache: cache[airly.Measurement]{
			d:   make(map[uint]entry[airly.Measurement], 50),
			ttl: airlyExpire,
		},
		depsCache: cache[[]typedDeparture]{
			d:   make(map[uint]entry[[]typedDeparture], 50),
			ttl: depsExpire,
		},
		airlyCli:    airly.NewClient(airly.URL),
		ttssBusCli:  ttss.NewClient(ttss.Bus),
		ttssTramCli: ttss.NewClient(ttss.Tram),
		searchCli:   s,
	}, nil
}
