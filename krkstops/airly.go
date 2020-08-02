package krkstops

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
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

type installation struct {
	ID int32
}

const airlyMeasurementsURL = "https://airapi.airly.eu/v2/measurements/installation"
const airlyInstallationsURL = "https://airapi.airly.eu/v2/installations/nearest"

// GetAirly queries external API and parses response
func (app *App) GetAirly(installation *pb.Installation) (pb.Airly, error) {
	airly := pb.Airly{}
	req, err := http.NewRequest("GET", airlyMeasurementsURL, nil)
	if err != nil {
		return airly, err
	}
	req.Header.Add("apikey", os.Getenv("AIRLY_KEY"))
	q := req.URL.Query()
	q.Add("installationId", fmt.Sprint(installation.GetId()))
	req.URL.RawQuery = q.Encode()
	resp, err := app.HTTPClient.Do(req)
	if err != nil {
		return airly, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return airly, err
	}
	if resp.StatusCode != 200 {
		return airly, errors.New(string(body))
	}
	var airlyResp AirlyResp
	err = json.Unmarshal(body, &airlyResp)
	if err != nil {
		return airly, err
	}
	for _, index := range airlyResp.Current.Indexes {
		if index.Name == "AIRLY_CAQI" {
			airly.Caqi = int32(index.Value)
			airly.Color = index.Color
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
	return airly, err
}

// GetAirlyInstallation queries Airly API for nearest installation
func (app *App) GetAirlyInstallation(position *pb.Positon) (*pb.Installation, error) {
	installation := make([]pb.Installation, 1)
	req, err := http.NewRequest("GET", airlyInstallationsURL, nil)
	if err != nil {
		return &installation[0], err
	}
	req.Header.Add("apikey", os.Getenv("AIRLY_KEY"))
	q := req.URL.Query()
	q.Add("lat", fmt.Sprint(position.Latitude))
	q.Add("lng", fmt.Sprint(position.Longitude))
	req.URL.RawQuery = q.Encode()
	resp, err := app.HTTPClient.Do(req)
	if err != nil {
		return &installation[0], err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &installation[0], err
	}
	if resp.StatusCode != 200 {
		return &installation[0], errors.New(string(body))
	}
	err = json.Unmarshal(body, &installation)
	return &installation[0], err
}
