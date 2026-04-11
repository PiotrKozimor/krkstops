package gtfs

import (
	"encoding/csv"
	"fmt"
	"strconv"
)

type Route struct {
	Name uint32
}

type Routes map[uint32]Route

func (u *Unmarshaler) UnmarshalRoutes(r *csv.Reader) (Routes, error) {
	routes := make(Routes, 100)
	r.Read()
	err := u.iterateSorted(r, 0, func(i int, record []string) error {
		id := uint32(i)
		u.routeIds[record[0]] = id
		var name int
		switch record[2] {
		case "LR0":
			name = 1000
		default:
			var err error
			name, err = strconv.Atoi(record[2])
			if err != nil {
				return fmt.Errorf("parse route name: %w", err)
			}
		}
		routes[id] = Route{
			Name: uint32(name),
		}
		return nil
	})
	return routes, err
}

func (u *Unmarshaler) routeId(s string) (uint32, error) {
	id, ok := u.routeIds[s]
	if !ok {
		return 0, fmt.Errorf("not found: %s", s)
	}
	return id, nil
}
