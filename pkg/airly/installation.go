package airly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	nearestInstallationsPath = "v2/installations/nearest"
	installationsPath        = "v2/installations/%d"
)

type installation struct {
	Id       int32
	Location struct {
		Latitude  float32
		Longitude float32
	}
}

// NearestInstallation for given location.
func (c Client) NearestInstallation(lat, lon float32) (Installation, error) {
	var i Installation

	r, close, err := c.do(nearestInstallationsPath, func(r *http.Request) {
		q := url.Values{}
		q.Add("lat", fmt.Sprint(lat))
		q.Add("lng", fmt.Sprint(lon))
		r.URL.RawQuery = q.Encode()
	})
	if err != nil {
		return i, fmt.Errorf("nearest installation: %w", err)
	}
	defer close()

	airlyInstallations := make([]installation, 1)
	enc := json.NewDecoder(r)
	err = enc.Decode(&airlyInstallations)
	if err != nil {
		return i, err
	}

	i.Id = airlyInstallations[0].Id
	i.Latitude = airlyInstallations[0].Location.Latitude
	i.Longitude = airlyInstallations[0].Location.Longitude
	return i, nil
}

// GetInstallation returns full details about installation with given ID
func (c Client) GetInstallation(id uint) (Installation, error) {
	var i Installation

	r, close, err := c.do(fmt.Sprintf(installationsPath, id), nil)
	if err != nil {
		return i, fmt.Errorf("get installation: %w", err)
	}
	defer close()

	var airlyInstallation installation
	enc := json.NewDecoder(r)
	err = enc.Decode(&airlyInstallation)
	if err != nil {
		return i, err
	}

	i.Id = airlyInstallation.Id
	i.Longitude = airlyInstallation.Location.Longitude
	i.Latitude = airlyInstallation.Location.Latitude
	return i, nil
}
