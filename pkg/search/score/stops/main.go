package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"strings"

	"github.com/PiotrKozimor/krkstops/pkg/stops"
	"github.com/PiotrKozimor/krkstops/pkg/ttss"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//go:generate go run .
func main() {
	cli := ttss.NewClient(ttss.Bus)
	allStops, err := cli.GetAllStops()
	handle(err)
	for i := range allStops {
		allStops[i].Name = trim(allStops[i].Name)
	}

	cliTram := ttss.NewClient(ttss.Tram)
	stopsTram, err := cliTram.GetAllStops()
	handle(err)
	for i := range stopsTram {
		stopsTram[i].Name = trim(stopsTram[i].Name)
	}

	stopsMerged, err := stops.Do(allStops, stopsTram)
	handle(err)

	buf := bytes.Buffer{}
	err = gob.NewEncoder(&buf).Encode(stopsMerged)
	handle(err)
	err = os.WriteFile("stops.gob", buf.Bytes(), 0644)
	handle(err)
}

func trim(name string) string {
	return strings.TrimSpace(strings.TrimSuffix(name, "(nż)"))
}
