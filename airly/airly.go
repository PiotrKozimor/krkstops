package airly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/PiotrKozimor/krkstops/pb"
)

type index struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Color string  `json:"color"`
}

type indicator struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type airlyParameters struct {
	Values       []indicator `json:"values"`
	Indexes      []index     `json:"indexes"`
	FromDateTime string      `json:"fromDateTime"`
}

// AirlyResp describes response from installation
type AirlyResp struct {
	Current airlyParameters `json:"current"`
}

type location struct {
	Latitude  float32
	Longitude float32
}
type installation struct {
	Id       int32
	Location location
}

const airlyMeasurementsURL = "https://airapi.airly.eu/v2/measurements/installation"
const airlyNearestInstallationsURL = "https://airapi.airly.eu/v2/installations/nearest"
const airlyInstallationsURL = "https://airapi.airly.eu/v2/installations/"

// GetAirly queries external API and parses response
func GetAirly(installation *pb.Installation) (*pb.Airly, error) {
	airly := pb.Airly{}
	req, err := http.NewRequest("GET", airlyMeasurementsURL, nil)
	if err != nil {
		return &airly, err
	}
	req.Header.Add("apikey", os.Getenv("AIRLY_KEY"))
	q := req.URL.Query()
	q.Add("installationId", fmt.Sprint(installation.GetId()))
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &airly, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return &airly, fmt.Errorf("get airly status code: %d", resp.StatusCode)
	}
	var airlyResp AirlyResp
	enc := json.NewDecoder(resp.Body)
	err = enc.Decode(&airlyResp)
	if err != nil {
		return &airly, err
	}
	for _, index := range airlyResp.Current.Indexes {
		if index.Name == "AIRLY_CAQI" {
			airly.Caqi = int32(index.Value)
			var color int64
			color, err = strconv.ParseInt(index.Color[1:], 16, 32)
			airly.Color = uint32(color)
		}
	}
	for _, indicator := range airlyResp.Current.Values {
		switch indicator.Name {
		case "HUMIDITY":
			airly.Humidity = int32(indicator.Value)
		case "TEMPERATURE":
			airly.Temperature = float32(indicator.Value)
		}
	}
	return &airly, err
}

// NearestInstallation for given location.
func NearestInstallation(location *pb.InstallationLocation) (*pb.Installation, error) {
	airlyInstallation := make([]installation, 1)
	var installation pb.Installation
	req, err := http.NewRequest("GET", airlyNearestInstallationsURL, nil)
	if err != nil {
		return &installation, err
	}
	req.Header.Add("apikey", os.Getenv("AIRLY_KEY"))
	q := req.URL.Query()
	q.Add("lat", fmt.Sprint(location.Latitude))
	q.Add("lng", fmt.Sprint(location.Longitude))
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &installation, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return &installation, fmt.Errorf("nearest installation status code: %d", resp.StatusCode)
	}
	enc := json.NewDecoder(resp.Body)
	err = enc.Decode(&airlyInstallation)
	return &pb.Installation{
		Id:        airlyInstallation[0].Id,
		Latitude:  airlyInstallation[0].Location.Latitude,
		Longitude: airlyInstallation[0].Location.Longitude,
	}, err
}

// GetInstallation returns full details about installation with given ID
func GetInstallation(id uint) (*pb.Installation, error) {
	var airlyInstallation installation
	req, err := http.NewRequest("GET", airlyInstallationsURL+fmt.Sprint(id), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("apikey", os.Getenv("AIRLY_KEY"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("get installation status code: %d", resp.StatusCode)
	}
	enc := json.NewDecoder(resp.Body)
	err = enc.Decode(&airlyInstallation)
	return &pb.Installation{
		Id:        airlyInstallation.Id,
		Longitude: airlyInstallation.Location.Longitude,
		Latitude:  airlyInstallation.Location.Latitude,
	}, err
}
