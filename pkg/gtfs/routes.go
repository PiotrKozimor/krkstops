package gtfs

import (
	"encoding/csv"
)

type Routes map[string]string

func (u *Unmarshaler) UnmarshalRoutes(r *csv.Reader) (Routes, error) {
	routes := make(Routes, 100)
	r.Read()
	err := u.iterate(r, func(record []string) error {
		routes[record[0]] = record[2]
		return nil
	})
	return routes, err
}
