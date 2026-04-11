package gtfs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

var (
	Bus = &Unmarshaler{
		estimatedTripSize: 50000,
	}
	Tram = &Unmarshaler{
		estimatedTripSize: 20000,
	}
	Mobilis = &Unmarshaler{
		estimatedTripSize: 30000,
	}
)

type Unmarshaler struct {
	estimatedTripSize int
	serviceIds        map[string]uint32
	tripIds           map[string]uint32
	stopIds           map[string]uint32
	routeIds          map[string]uint32
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
