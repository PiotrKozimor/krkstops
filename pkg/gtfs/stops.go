package gtfs

import (
	"encoding/csv"
	"fmt"
	"strings"
)

type Stops map[string][]uint32

func (u *Unmarshaler) UnmarshalStops(r *csv.Reader) (Stops, error) {
	stops := make(Stops, 500)
	r.Read()
	err := u.iterateSorted(r, 0, func(i int, record []string) (err error) {
		name := strings.TrimSuffix(record[2], " (nż)")
		id := uint32(i)
		u.stopIds[record[0]] = id
		stops[name] = append(stops[name], id)
		return nil
	})
	return stops, err
}

func (u *Unmarshaler) stopId(s string) (uint32, error) {
	id, ok := u.stopIds[s]
	if !ok {
		return 0, fmt.Errorf("not found: %s", s)
	}
	return id, nil
}
