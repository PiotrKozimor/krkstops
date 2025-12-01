package krkstops

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PiotrKozimor/krkstops/pkg/gtfs"
	"github.com/PiotrKozimor/krkstops/pkg/trie"
)

type scoredStop struct {
	name  string
	score uint32
}

func (s *KrkStopsServer) initGtfs() error {

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

		departure.GetFeed = func() ([]byte, error) {
			resp, err := http.Get(instance.updatesUrl)
			if err != nil {
				return nil, fmt.Errorf("http get: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				return nil, fmt.Errorf("status code: %d", resp.StatusCode)
			}
			return io.ReadAll(resp.Body)
		}
		err = departure.Init(func(file string) (*csv.Reader, error) {
			f, err := r.Open(file)
			return csv.NewReader(f), err
		})
		if err != nil {
			return fmt.Errorf("init %s: %w", instance.name, err)
		}
		s.departures = append(s.departures, departure)

		go departure.RefreshUpdates()
	}

	s.search = trie.New(func(s1, s2 scoredStop) int {
		return int(s2.score) - int(s1.score)
	}, func(s1, s2 scoredStop) bool {
		return s1.name == s2.name
	})
	unique := make(map[string]uint32)
	for _, d := range s.departures {
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
	s.search.InsertWords(entries...)

	return nil

}
