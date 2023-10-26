package search

import (
	"bytes"
	_ "embed"
	"encoding/gob"
	"slices"

	"github.com/PiotrKozimor/krkstops/pkg/search/merged"
	"github.com/PiotrKozimor/krkstops/pkg/trie"
)

var (
	//go:embed score/stops/stops.gob
	stopsB []byte
	//go:embed score/score.gob
	scoreB []byte
)

type Search struct {
	t     trie.Trie
	score map[uint]uint
	stops merged.Stops
}

type Stop struct {
	Name      string
	Id        uint
	Score     uint
	Tram, Bus bool
}

type scoredStop struct {
	id, score uint
}

func New() (*Search, error) {
	stops, err := merged.Read(stopsB)
	if err != nil {
		return nil, err
	}
	s := Search{
		t:     trie.New(),
		stops: stops,
	}
	buf := bytes.NewBuffer(scoreB)
	err = gob.NewDecoder(buf).Decode(&s.score)
	if err != nil {
		return nil, err
	}
	entries := make([]trie.Entry, 0, len(stops))
	for id, stop := range stops {
		entries = append(entries, trie.Entry{
			Id:   id,
			Word: stop.Name,
		})
	}
	s.t.InsertWords(entries...)
	return &s, nil
}

func (s *Search) Search(term string, limit int) []Stop {
	r := s.t.SearchExact(term)
	stops := s.sort(r, limit)
	if len(stops) >= limit {
		return stops
	}
	extra := s.t.SearchInDistance(term, 1)
	extraStops := s.sort(extra, limit-len(stops))
	stops = append(stops, extraStops...)
	return stops
}

func (s *Search) searchWithinDistance(term string, limit int) []Stop {
	r := s.t.SearchWithinDistance(term, 1)
	return s.sort(r, limit)
}

func (s *Search) SearchExact(term string, limit int) []Stop {
	r := s.t.SearchExact(term)
	return s.sort(r, limit)
}

func (s *Search) sort(results []uint, limit int) []Stop {
	scored := make([]scoredStop, len(results))
	for i := range results {
		scored[i] = scoredStop{id: results[i], score: s.score[results[i]]}
	}
	slices.SortFunc(scored, func(a, b scoredStop) int {
		return int(a.score - b.score)
	})
	res := make([]Stop, 0, limit)
	for i, sc := range scored {
		if i >= limit {
			break
		}
		res = append(res, Stop{
			Name:  s.stops[sc.id].Name,
			Bus:   s.stops[sc.id].Bus,
			Tram:  s.stops[sc.id].Tram,
			Score: sc.score,
			Id:    sc.id,
		})
	}
	return res
}

func (s *Search) Get(id uint) Stop {
	stop := s.stops[id]
	return Stop{
		Name:  stop.Name,
		Bus:   stop.Bus,
		Tram:  stop.Tram,
		Score: s.score[id],
		Id:    id,
	}
}
