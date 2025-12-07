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
	err := u.iterate(r, func(record []string) error {
		id, err := u.reduceRouteId(record[0])
		if err != nil {
			return fmt.Errorf("reduce route id: %w", err)
		}
		name, err := strconv.Atoi(record[2])
		routes[id] = Route{
			Name: uint32(name),
		}
		return err
	})
	return routes, err
}
