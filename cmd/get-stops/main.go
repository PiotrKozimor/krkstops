package main

import (
	"flag"
	"log"
	"strings"

	"go.cozymore.dev/krkstops/pkg/store"
	"go.cozymore.dev/krkstops/pkg/ttss"
)

const (
	bus  = "bus"
	tram = "tram"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	legalNoticeTemplatePath := flag.String("legal", "", "Tool will read legal_notice.tmpl.txt file from this path and write it back to the same folder as legal_notice.txt")
	transit := flag.String("transit", "", "[bus, tram] Tool will retrieve bus or tram stops only, both if empty")
	flag.Parse()

	var cli *ttss.Client
	switch *transit {
	case "":
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

		merged, err := store.Merge(allStops, stopsTram)
		handle(err)

		err = store.Write(merged, *legalNoticeTemplatePath)
		handle(err)
		return
	case bus:
		cli = ttss.NewClient(ttss.Bus)
	case tram:
		cli = ttss.NewClient(ttss.Tram)
	default:
		log.Fatal("unrecognised transit flag: ", *transit)
	}
	allStops, err := cli.GetAllStops()
	handle(err)

	stopsToStore := make(store.Stops, len(allStops))
	for i := range allStops {
		stopsToStore[allStops[i].Id] = store.Stop{
			Name: trim(allStops[i].Name),
			Tram: *transit == tram,
			Bus:  *transit == bus,
		}
		allStops[i].Name = trim(allStops[i].Name)
	}

	err = store.Write(stopsToStore, *legalNoticeTemplatePath)
	handle(err)
}

func trim(name string) string {
	return strings.TrimSpace(strings.TrimSuffix(name, "(nż)"))
}
