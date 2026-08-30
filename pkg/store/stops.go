package store

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"time"

	"go.cozymore.dev/krkstops/pkg/ttss"
)

const (
	stopsFile    = "stops.gob"
	templateFile = "legal_notice.tmpl.txt"
	noticeFile   = "legal_notice.txt"
)

type Stops map[uint]Stop

type Stop struct {
	Name      string
	Tram, Bus bool
}

func Merge(bus, tram []ttss.Stop) (Stops, error) {
	stopsM := make(Stops, len(bus))
	for _, stop := range bus {
		stopsM[stop.Id] = Stop{Name: stop.Name, Bus: true}
	}

	for _, stop := range tram {
		if s, ok := stopsM[stop.Id]; ok {
			if s.Name != stop.Name {
				return nil, fmt.Errorf("got different name for stop %d\tbus: '%s'\ttram: '%s'", stop.Id, s.Name, stop.Name)
			}
			s.Tram = true
			stopsM[stop.Id] = s
		} else {
			stopsM[stop.Id] = Stop{
				Name: stop.Name,
				Tram: true,
			}
		}
	}

	return stopsM, nil
}

func Read(b []byte) (Stops, error) {
	var s Stops
	buf := bytes.NewBuffer(b)
	err := gob.NewDecoder(buf).Decode(&s)
	return s, err
}

func ReadDefault() (Stops, error) {
	b, err := os.ReadFile(stopsFile)
	if err != nil {
		return nil, err
	}
	return Read(b)
}

func Write(merged Stops, legalNoticeTemplatePath string) error {
	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(merged)
	if err != nil {
		return err
	}
	err = os.WriteFile(stopsFile, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	b, err := os.ReadFile(path.Join(legalNoticeTemplatePath, templateFile))
	if err != nil {
		return err
	}
	notice := fmt.Sprintf(string(b), time.Now().Format(time.DateTime))
	err = os.WriteFile(path.Join(legalNoticeTemplatePath, noticeFile), []byte(notice), 0644)
	return err
}
