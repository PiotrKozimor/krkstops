package krkstops

import (
	"log"
	"sync"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/pkg/trie"
	"github.com/PiotrKozimor/krkstops/pkg/ttss"
	"github.com/PiotrKozimor/krkstops/pkg/ttssstops"
)

type KrkStopsServer struct {
	pb.UnimplementedKrkStopsServer
	depsCache     cache[[]typedDeparture]
	ttssBusCli    ttssClient
	ttssTramCli   ttssClient
	searchCli     *ttssstops.Search
	searchMu      sync.RWMutex
	search        trie.Trie[scoredStop]
	departures    []*gtfsDepartures
	departuresMu  sync.RWMutex
	updatesCancel func()
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
		ttssBusCli:    ttss.NewClient(ttss.Bus),
		ttssTramCli:   ttss.NewClient(ttss.Tram),
		searchCli:     s,
		updatesCancel: func() {},
	}
	now := time.Now()
	err = server.initGtfs()
	log.Print("init gtfs took: ", time.Since(now))
	go server.refreshGtfs()

	return server, err
}
