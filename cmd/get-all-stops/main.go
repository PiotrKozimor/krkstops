package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/PiotrKozimor/krkstops/pkg/stops"
	"github.com/PiotrKozimor/krkstops/pkg/ttss"
)

const (
	templateFile = "legal_notice.tmpl.txt"
	noticeFile   = "legal_notice.txt"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	stopsFile := flag.String("stops", "stops.gob", "Tool will write stops file to this path")
	legalNoticeTemplatePath := flag.String("legal", "", "Tool will read legal_notice.tmpl.txt file from this path and write it back to the same folder as legal_notice.txt")
	flag.Parse()

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

	stopsMerged, err := stops.Merge(allStops, stopsTram)
	handle(err)

	buf := bytes.Buffer{}
	err = gob.NewEncoder(&buf).Encode(stopsMerged)
	handle(err)
	err = os.WriteFile(*stopsFile, buf.Bytes(), 0644)
	handle(err)

	b, err := os.ReadFile(path.Join(*legalNoticeTemplatePath, templateFile))
	handle(err)
	notice := fmt.Sprintf(string(b), time.Now().Format(time.DateTime))
	err = os.WriteFile(path.Join(*legalNoticeTemplatePath, noticeFile), []byte(notice), 0644)
	handle(err)
}

func trim(name string) string {
	return strings.TrimSpace(strings.TrimSuffix(name, "(nż)"))
}
