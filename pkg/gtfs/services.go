package gtfs

import (
	"encoding/csv"
	"fmt"
	"slices"
	"time"
)

const (
	layout = "20060102"
)

type Service struct {
	Id       uint32
	Days     []time.Weekday
	StartsAt time.Time
	EndsAt   time.Time
}

type ServiceException struct {
	ServiceId uint32
	StartsAt  time.Time
}

func (u *Unmarshaler) UnmarshalServices(r *csv.Reader) ([]Service, error) {
	var services []Service
	r.Read()
	err := u.iterate(r, func(record []string) (err error) {
		service := Service{}
		service.Id, err = u.reduceServiceId(record[0])
		if err != nil {
			return fmt.Errorf("reduce service id: %w", err)
		}
		for i := range 7 {
			if record[1+i] == "1" {
				weekday := i + 1
				if weekday >= 7 {
					weekday = 0
				}
				service.Days = append(service.Days, time.Weekday(weekday))
			}
		}
		service.StartsAt, err = time.Parse(layout, record[8])
		if err != nil {
			return fmt.Errorf("parse start at: %w", err)
		}
		service.EndsAt, err = time.Parse(layout, record[9])
		if err != nil {
			return fmt.Errorf("parse end at: %w", err)
		}
		services = append(services, service)
		return nil
	})
	return services, err
}

func (u *Unmarshaler) UnmarshalServiceExceptions(r *csv.Reader) ([]ServiceException, error) {
	var exceptions []ServiceException
	r.Read()
	err := u.iterate(r, func(record []string) (err error) {
		exception := ServiceException{}
		exception.ServiceId, err = u.reduceServiceId(record[0])
		if err != nil {
			return fmt.Errorf("reduce service id: %w", err)
		}
		exception.StartsAt, err = time.Parse(layout, record[1])
		if err != nil {
			return fmt.Errorf("parse start at: %w", err)
		}
		exceptions = append(exceptions, exception)
		return nil
	})
	slices.SortFunc(exceptions, func(a, b ServiceException) int {
		return a.StartsAt.Compare(b.StartsAt)
	})
	return exceptions, err
}
