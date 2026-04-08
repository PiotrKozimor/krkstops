package gtfs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

var (
	Bus = &Unmarshaler{
		estimatedTripSize: 50000,
		reduceRouteId: func(s string) (uint32, error) {
			switch s {
			case "LR0":
				return 990, nil
			default:
				id, err := strconv.Atoi(s)
				return uint32(id), err
			}
		},
	}
	Tram = &Unmarshaler{
		estimatedTripSize: 20000,
		reduceRouteId: func(s string) (uint32, error) {
			trim := strings.TrimPrefix(s, "route_")
			id, err := strconv.Atoi(trim)
			return uint32(id), err
		},
	}
	Mobilis = &Unmarshaler{
		estimatedTripSize: 30000,
		reduceRouteId: func(s string) (uint32, error) {
			id, err := strconv.Atoi(s)
			return uint32(id), err
		},
	}
)

type Unmarshaler struct {
	estimatedTripSize int
	serviceIds        map[string]uint32
	tripIds           map[string]uint32
	stopIds           map[string]uint32
	reduceRouteId     func(string) (uint32, error)
}

func (u *Unmarshaler) iterate(r *csv.Reader, c func([]string) error) error {
	for {
		record, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		err = c(record)
		if err != nil {
			return err
		}
	}
}

func (u *Unmarshaler) iterateSorted(r *csv.Reader, compareAtColumn int, c func(int, []string) error) error {
	allRecords, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}
	slices.SortFunc(allRecords, func(a, b []string) int {
		return strings.Compare(a[compareAtColumn], b[compareAtColumn])
	})
	for i, record := range allRecords {
		err = c(i, record)
		if err != nil {
			return err
		}
	}
	return nil
}
