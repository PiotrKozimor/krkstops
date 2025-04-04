package krkstops

import (
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/pkg/search"
	"github.com/PiotrKozimor/krkstops/pkg/ttss"
)

type KrkStopsServer struct {
	pb.UnimplementedKrkStopsServer
	depsCache   cache[[]typedDeparture]
	ttssBusCli  ttssClient
	ttssTramCli ttssClient
	searchCli   *search.Search
}

type ttssClient interface {
	GetDepartures(uint) ([]ttss.Departure, error)
}

func NewServer() (*KrkStopsServer, error) {
	s, err := search.New()
	if err != nil {
		return nil, err
	}
	return &KrkStopsServer{
		depsCache: cache[[]typedDeparture]{
			d:   make(map[uint]entry[[]typedDeparture], 50),
			ttl: depsExpire,
		},
		ttssBusCli:  ttss.NewClient(ttss.Bus),
		ttssTramCli: ttss.NewClient(ttss.Tram),
		searchCli:   s,
	}, nil
}
