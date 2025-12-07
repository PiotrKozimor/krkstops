package krkstops

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PiotrKozimor/krkstops/pkg/gtfs"
	"github.com/PiotrKozimor/krkstops/pkg/trie"
	"github.com/prometheus/client_golang/prometheus"
)

type scoredStop struct {
	name  string
	score uint32
}

type gtfsDepartures struct {
	*gtfs.Departures
	getUpdates func() (gtfs.StopUpdates, error)
}

func (s *KrkStopsServer) initGtfs() error {
	departures := make([]*gtfsDepartures, 0, 3)

	for _, instance := range []struct {
		name       string
		u          *gtfs.Unmarshaler
		url        string
		updatesUrl string
	}{
		{"bus", gtfs.Bus, "https://gtfs.ztp.krakow.pl/GTFS_KRK_A.zip", "https://gtfs.ztp.krakow.pl/TripUpdates_A.pb"},
		{"mobilis", gtfs.Mobilis, "https://gtfs.ztp.krakow.pl/GTFS_KRK_M.zip", "https://gtfs.ztp.krakow.pl/TripUpdates_M.pb"},
		{"tram", gtfs.Tram, "https://gtfs.ztp.krakow.pl/GTFS_KRK_T.zip", "https://gtfs.ztp.krakow.pl/TripUpdates_T.pb"},
	} {
		resp, err := http.Get(instance.url)
		if err != nil {
			return fmt.Errorf("get gtfs: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("status code: %d", resp.StatusCode)
		}

		b := &bytes.Buffer{}
		_, err = io.Copy(b, resp.Body)
		if err != nil {
			return fmt.Errorf("copy response: %w", err)
		}
		r, err := zip.NewReader(bytes.NewReader(b.Bytes()), int64(b.Len()))
		if err != nil {
			return fmt.Errorf("zip reader: %w", err)
		}
		departure := gtfs.NewDepartures(instance.u)

		err = departure.Init(func(file string) (*csv.Reader, error) {
			f, err := r.Open(file)
			return csv.NewReader(f), err
		})
		if err != nil {
			return fmt.Errorf("init %s: %w", instance.name, err)
		}
		departures = append(departures, &gtfsDepartures{
			Departures: departure,
			getUpdates: fetchUpdates(instance.u, instance.updatesUrl),
		})
	}

	search := trie.New(func(s1, s2 scoredStop) int {
		return int(s2.score) - int(s1.score)
	}, func(s1, s2 scoredStop) bool {
		return s1.name == s2.name
	})

	unique := make(map[string]uint32)
	for _, d := range departures {
		for name, ids := range d.Stops {
			name = strings.Replace(name, "”", "\"", -1)
			for _, id := range ids {
				unique[name] += d.StopsScore[id]
			}
		}
	}
	entries := make([]trie.Entry[scoredStop], 0, 5000)
	for name, score := range unique {
		entries = append(entries, trie.Entry[scoredStop]{
			Item: scoredStop{
				name:  name,
				score: score,
			},
			Word: name,
		})

	}
	search.InsertWords(entries...)

	s.updatesCancel()
	ctx, cancel := context.WithCancel(context.Background())
	for _, d := range departures {
		go refreshUpdates(d, ctx)
	}
	s.updatesCancel = cancel

	s.searchMu.Lock()
	s.search = search
	s.searchMu.Unlock()

	s.departuresMu.Lock()
	s.departures = departures
	s.departuresMu.Unlock()

	return nil
}

func (s *KrkStopsServer) refreshGtfs() {
	ticker := time.NewTicker(time.Hour)
	for t := range ticker.C {
		if t.In(gtfs.Location).Hour() == 3 {
			err := s.initGtfs()
			if err != nil {
				serverErrors.With(prometheus.Labels{"err": "refresh_gtfs"}).Inc()
				log.Print("failed to refresh gtfs data: ", err)
			} else {
				serverSuccesses.With(prometheus.Labels{"name": "refresh_gtfs"}).Inc()
				log.Print("gtfs data refreshed")
			}
		}
	}
}

func fetchUpdates(u *gtfs.Unmarshaler, url string) func() (gtfs.StopUpdates, error) {
	return func() (gtfs.StopUpdates, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("http get: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("status code: %d", resp.StatusCode)
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read all: %w", err)
		}
		feed, err := gtfs.ParseFeed(b)
		if err != nil {
			return nil, fmt.Errorf("parse feed: %w", err)
		}
		return u.ParseStopUpdates(feed)
	}
}

func refreshUpdates(d *gtfsDepartures, ctx context.Context) {
	log.Print("starting refresh goroutine")
	defer log.Print("terminating refresh goroutine")
	t := time.NewTicker(time.Second * 20)
	for {
		updates, err := d.getUpdates()
		if err != nil {
			serverErrors.With(prometheus.Labels{"err": "get_updates"}).Inc()
			log.Print("failed to get updates", err)
		} else {
			serverSuccesses.With(prometheus.Labels{"name": "get_updates"}).Inc()
		}
		d.SetStopUpdates(updates)
		select {
		case <-t.C:
		case <-ctx.Done():
			return
		}
	}
}
