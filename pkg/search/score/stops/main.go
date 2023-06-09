package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"strings"

	"github.com/PiotrKozimor/krkstops/pkg/search/merged"
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
	stops, err := cli.GetAllStops()
	handle(err)
	for i := range stops {
		stops[i].Name = trim(stops[i].Name)
	}

	cliTram := ttss.NewClient(ttss.Tram)
	stopsTram, err := cliTram.GetAllStops()
	handle(err)
	for i := range stopsTram {
		stopsTram[i].Name = trim(stopsTram[i].Name)
	}

	stopsMerged, err := merged.Do(stops, stopsTram)
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
