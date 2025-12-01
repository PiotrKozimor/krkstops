package krkstops

import (
	"log"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/pkg/gtfs"
	"github.com/PiotrKozimor/krkstops/pkg/trie"
	"github.com/PiotrKozimor/krkstops/pkg/ttss"
	"github.com/PiotrKozimor/krkstops/pkg/ttssstops"
)

type KrkStopsServer struct {
	pb.UnimplementedKrkStopsServer
	depsCache   cache[[]typedDeparture]
	ttssBusCli  ttssClient
	ttssTramCli ttssClient
	searchCli   *ttssstops.Search
	search      trie.Trie[scoredStop]
	departures  []*gtfs.Departures
}

type ttssClient interface {
	GetDepartures(uint) ([]ttss.Departure, error)
}

func NewServer() (*KrkStopsServer, error) {
	s, err := ttssstops.New()
	if err != nil {
		return nil, err
	}
	server := &KrkStopsServer{
		depsCache: cache[[]typedDeparture]{
			d:   make(map[uint]entry[[]typedDeparture], 50),
			ttl: depsExpire,
		},
		ttssBusCli:  ttss.NewClient(ttss.Bus),
		ttssTramCli: ttss.NewClient(ttss.Tram),
		searchCli:   s,
	}
	now := time.Now()
	err = server.initGtfs()
	log.Print("init gtfs took: ", time.Since(now))

	return server, err
}
