package ttssstops

import (
	"bytes"
	_ "embed"
	"encoding/gob"

	"github.com/PiotrKozimor/krkstops/pkg/store"
	"github.com/PiotrKozimor/krkstops/pkg/trie"
)

//go:generate go run ../../cmd/get-stops -legal ../../cmd/krkstops
//go:generate go run ../../cmd/score-stops
var (
	//go:embed stops.gob
	stopsB []byte
	//go:embed score.gob
	scoreB []byte
)

type Search struct {
	trie.Trie[Stop]
	score map[uint]uint
	stops store.Stops
}

type Stop struct {
	Name      string
	Id        uint
	score     uint
	Tram, Bus bool
}

func New() (*Search, error) {
	stops, err := store.Read(stopsB)
	if err != nil {
		return nil, err
	}
	s := Search{
		stops: stops,
	}
	s.Trie = trie.New(s.compare, s.equal)

	buf := bytes.NewBuffer(scoreB)
	err = gob.NewDecoder(buf).Decode(&s.score)
	if err != nil {
		return nil, err
	}
	entries := make([]trie.Entry[Stop], 0, len(stops))
	for id, stop := range stops {
		entries = append(entries, trie.Entry[Stop]{
			Item: Stop{
				score: s.score[id],
				Name:  stop.Name,
				Id:    id,
				Tram:  stop.Tram,
				Bus:   stop.Bus,
			},
			Word: stop.Name,
		})
	}
	s.InsertWords(entries...)
	return &s, nil
}

func (s *Search) compare(a, b Stop) int {
	return int(b.score) - int(a.score)
}

func (s *Search) equal(a, b Stop) bool {
	return a.score == b.score
}

func (s *Search) Get(id uint) Stop {
	stop := s.stops[id]
	return Stop{
		Name:  stop.Name,
		Bus:   stop.Bus,
		Tram:  stop.Tram,
		score: s.score[id],
		Id:    id,
	}
}
