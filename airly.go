package krkstops

import (
	"context"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/pkg/airly"
)

func (s *KrkStopsServer) GetAirly(ctx context.Context, req *pb.GetMeasurementRequest) (*pb.Measurement, error) {
	var err error
	measurement, _, ok := s.airlyCache.get(uint(req.Id))
	if !ok {
		measurement, err = s.airlyCli.GetMeasurement(uint(req.Id))
		if err != nil {
			return nil, err
		}
		s.airlyCache.set(uint(req.Id), measurement)
	}

	return protoMeasurement(measurement), err
}

func (s *KrkStopsServer) FindNearestAirlyInstallation(ctx context.Context, location *pb.Location) (*pb.Installation, error) {
	inst, err := s.airlyCli.NearestInstallation(location.Latitude, location.Longitude)
	return protoInstallation(inst), err
}

func (s *KrkStopsServer) GetAirlyInstallation(ctx context.Context, req *pb.GetAirlyInstallationRequest) (*pb.Installation, error) {
	inst, err := s.airlyCli.GetInstallation(uint(req.Id))
	return protoInstallation(inst), err
}

func protoMeasurement(m airly.Measurement) *pb.Measurement {
	a := &pb.Measurement{
		Caqi:        m.Caqi,
		Humidity:    m.Humidity,
		Temperature: m.Temperature,
		Color:       m.Color,
	}
	return a
}

func protoInstallation(i airly.Installation) *pb.Installation {
	return &pb.Installation{
		Id:        i.Id,
		Latitude:  i.Latitude,
		Longitude: i.Latitude,
	}
}
