package airly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type measurementResp struct {
	Current struct {
		Values []struct {
			Name  string  `json:"name"`
			Value float64 `json:"value"`
		} `json:"values"`
		Indexes []struct {
			Name  string  `json:"name"`
			Value float64 `json:"value"`
			Color string  `json:"color"`
		} `json:"indexes"`
		FromDateTime string `json:"fromDateTime"`
	} `json:"current"`
}

type Endpoint string

const measurementsPath = "v2/measurements/installation"

// GetAirly queries external API and parses response
func (c Client) GetMeasurement(installationId uint) (Measurement, error) {
	var m Measurement

	r, close, err := c.do(measurementsPath, func(r *http.Request) {
		q := url.Values{}
		q.Add("installationId", fmt.Sprint(installationId))
		r.URL.RawQuery = q.Encode()
	})
	if err != nil {
		return m, fmt.Errorf("get measurement: %w", err)
	}
	defer close()

	var airlyMeasurement measurementResp
	enc := json.NewDecoder(r)
	err = enc.Decode(&airlyMeasurement)
	if err != nil {
		return m, err
	}

	for _, index := range airlyMeasurement.Current.Indexes {
		if index.Name == "AIRLY_CAQI" {
			m.Caqi = int32(index.Value)
			var color int64
			color, err = strconv.ParseInt(index.Color[1:], 16, 32)
			if err != nil {
				return m, err
			}
			m.Color = uint32(color)
		}
	}
	for _, indicator := range airlyMeasurement.Current.Values {
		switch indicator.Name {
		case "HUMIDITY":
			m.Humidity = int32(indicator.Value)
		case "TEMPERATURE":
			m.Temperature = float32(indicator.Value)
		}
	}
	return m, nil
}
