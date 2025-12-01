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
	err := u.iterate(r, func(record []string) (err error) {
		id, err := u.reduceStopId(record[0])
		if err != nil {
			return fmt.Errorf("reduce stop id: %w", err)
		}
		name := strings.TrimSuffix(record[2], " (nż)")
		stops[name] = append(stops[name], id)
		return nil
	})
	return stops, err
}
