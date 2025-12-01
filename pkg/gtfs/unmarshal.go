package gtfs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

var (
	reduceServiceId = func(s string) (uint32, error) {
		trim := strings.TrimPrefix(s, "service_")
		id, err := strconv.Atoi(trim)
		return uint32(id), err
	}
	reduceRouteId = func(s string) (uint32, error) {
		trim := strings.TrimPrefix(s, "route_")
		id, err := strconv.Atoi(trim)
		return uint32(id), err
	}
	reduceTripId = func(s string) (uint32, error) {
		split := strings.Split(s, "_")
		if len(split) != 6 {
			return 0, fmt.Errorf("invalid id: %s", s)
		}
		id1, err1 := strconv.Atoi(split[1])
		id2, err2 := strconv.Atoi(split[3])
		id3, err3 := strconv.Atoi(split[5])
		if id1 > math.MaxUint16 {
			return 0, fmt.Errorf("invalid id1: %d", id1)
		}
		if id2 > math.MaxUint8 {
			return 0, fmt.Errorf("invalid id2: %d", id1)
		}
		if id3 > math.MaxUint8 {
			return 0, fmt.Errorf("invalid id3: %d", id1)
		}
		return uint32(id1) | uint32(id2)<<16 | uint32(id3)<<24, errors.Join(err1, err2, err3)
	}
	reduceStopId = func(s string) (uint32, error) {
		split := strings.Split(s, "_")
		if len(split) != 3 {
			return 0, fmt.Errorf("invalid id: %s", s)
		}
		id, err := strconv.Atoi(split[2])
		return uint32(id), err
	}

	Bus = &Unmarshaler{
		estimatedTripSize: 50000,
		reduceServiceId:   reduceServiceId,
		reduceRouteId:     reduceRouteId,
		reduceTripId:      reduceTripId,
		reduceStopId:      reduceStopId,
	}
	Tram = &Unmarshaler{
		estimatedTripSize: 20000,
		reduceServiceId:   reduceServiceId,
		reduceRouteId:     reduceRouteId,
		reduceTripId:      reduceTripId,
		reduceStopId:      reduceStopId,
	}
	Mobilis = &Unmarshaler{
		estimatedTripSize: 30000,
		reduceServiceId: func(s string) (uint32, error) {
			split := strings.Split(s, "_")
			if len(split) != 2 {
				return 0, fmt.Errorf("invalid id: %s", s)
			}
			id, err := strconv.Atoi(split[0])
			if len(split[1]) != 2 {
				return 0, fmt.Errorf("invalid suffix: %s", split[1])
			}
			serviceId := uint32(id) | uint32(split[1][0])<<16 | uint32(split[1][1])<<24
			return serviceId, err
		},
		reduceRouteId: func(s string) (uint32, error) {
			id, err := strconv.Atoi(s)
			return uint32(id), err
		},
		reduceTripId: func(s string) (uint32, error) {
			split := strings.Split(s, "_")
			id1, err1 := strconv.Atoi(split[0])
			id2, err2 := strconv.Atoi(split[1])
			if id1 > math.MaxUint16>>2 {
				return 0, fmt.Errorf("invalid id1: %d", id1)
			}
			if id2 > math.MaxUint16<<2 {
				return 0, fmt.Errorf("invalid id2: %d", id2)
			}
			return uint32(id1) | uint32(id2)<<14, errors.Join(err1, err2)
		},
		reduceStopId: func(s string) (uint32, error) {
			id, err := strconv.Atoi(s)
			return uint32(id), err
		},
	}
)

type Unmarshaler struct {
	estimatedTripSize int
	reduceServiceId   func(string) (uint32, error)
	reduceRouteId     func(string) (uint32, error)
	reduceTripId      func(string) (uint32, error)
	reduceStopId      func(string) (uint32, error)
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
