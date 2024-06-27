package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"time"

	_ "embed"

	"github.com/PiotrKozimor/krkstops/pkg/stops"
	"github.com/PiotrKozimor/krkstops/pkg/ttss"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	stopsB, err := os.ReadFile("stops.gob")
	handle(err)
	stops, err := stops.Read(stopsB)
	handle(err)

	scores := make(map[uint]uint, len(stops))

	cliB := ttss.NewClient(ttss.Bus)
	cliT := ttss.NewClient(ttss.Tram)
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
		time.Sleep(time.Millisecond * 100)
	}

	buf := bytes.Buffer{}
	err = gob.NewEncoder(&buf).Encode(scores)
	handle(err)
	err = os.WriteFile("score.gob", buf.Bytes(), 0644)
	handle(err)
}
