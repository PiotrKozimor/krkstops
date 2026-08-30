package krkstops

import (
	"log"
	"sync"
	"time"

	"go.cozymore.dev/krkstops/pkg/trie"
	"go.cozymore.dev/krkstops/pkg/ttss"
)

type KrkStopsServer struct {
	ttssBusCli    ttssClient
	ttssTramCli   ttssClient
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
	server := &KrkStopsServer{
		ttssBusCli:    ttss.NewClient(ttss.Bus),
		ttssTramCli:   ttss.NewClient(ttss.Tram),
		updatesCancel: func() {},
	}
	now := time.Now()
	err := server.initGtfs()
	log.Print("init gtfs took: ", time.Since(now))
	go server.refreshGtfs()

	return server, err
}
