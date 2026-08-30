package main

import (
	"log"
	"time"

	_ "embed"

	"go.cozymore.dev/krkstops/pkg/store"
	"go.cozymore.dev/krkstops/pkg/ttss"
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

	withRetries := func(c *ttss.Client, id uint) int {
		for range 3 {
			deps, err := c.GetDepartures(id)
			if err == nil {
				return len(deps)
			} else {
				time.Sleep(time.Millisecond * 500)
				continue
			}
		}
		log.Printf("failed to fetch score: %d", id)
		return 0
	}

	for id, stop := range stops {
		if stop.Bus {
			inc := withRetries(cliB, id)
			log.Printf("increasing score of stop with id %d by %d", id, inc)
			scores[id] += uint(inc)
		}
		if stop.Tram {
			inc := withRetries(cliT, id)
			log.Printf("increasing score of stop with id %d by %d", id, inc)
			scores[id] += uint(inc)
		}
		time.Sleep(time.Millisecond * 200)
	}

	err = store.WriteScore(scores)
	handle(err)
}
