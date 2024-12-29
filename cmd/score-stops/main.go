package main

import (
	"log"
	"time"

	_ "embed"

	"github.com/PiotrKozimor/krkstops/pkg/store"
	"github.com/PiotrKozimor/krkstops/pkg/ttss"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	stops, err := store.ReadDefault()
	handle(err)

	scores := make(map[uint]uint, len(stops))

	cliB := ttss.NewClient(ttss.Bus, ttss.WithTimeout(time.Second*5))
	cliT := ttss.NewClient(ttss.Tram, ttss.WithTimeout(time.Second*5))
	for id, stop := range stops {
		if stop.Bus {
			deps, err := cliB.GetDepartures(id)
			handle(err)
			inc := len(deps)
			log.Printf("increasing score of stop with id %d by %d", id, inc)
			scores[id] += uint(inc)
		}
		if stop.Tram {
			deps, err := cliT.GetDepartures(id)
			handle(err)
			inc := len(deps)
			log.Printf("increasing score of stop with id %d by %d", id, inc)
			scores[id] += uint(inc)
		}
		time.Sleep(time.Millisecond * 200)
	}

	err = store.WriteScore(scores)
	handle(err)
}
