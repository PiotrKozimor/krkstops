package merged

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/PiotrKozimor/krkstops/pkg/ttss"
)

type Stops map[uint]Stop

type Stop struct {
	Name      string
	Tram, Bus bool
}

func Do(bus, tram []ttss.Stop) (Stops, error) {
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
