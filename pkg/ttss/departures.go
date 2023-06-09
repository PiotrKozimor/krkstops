package ttss

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrStatusCode    = errors.New("got non 200 status code")
	ErrStopNotFound  = errors.New("stop not found")
	ErrRequestFailed = errors.New("request failed")
)

type departure struct {
	PlannedTime        string
	Status             string
	ActualRelativeTime int32
	Direction          string
	PatternText        string
}

type stopDepartures struct {
	Actual []departure
}

var (
	departuresPath = "internetservice/services/passageInfo/stopPassages/stop?stop=%d&mode=departure&language=pl"
)

// GetDepartures from Endpoint for stop with given shortName.
func (c Client) GetDepartures(id uint) ([]Departure, error) {
	url := fmt.Sprintf(strings.Join([]string{c.host, departuresPath}, "/"), id)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			return nil, fmt.Errorf("%w: %d", ErrStopNotFound, id)
		}
		return nil, fmt.Errorf("%w: %d", ErrStatusCode, resp.StatusCode)
	}
	var ttssDepartures stopDepartures
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ttssDepartures)
	if err != nil {
		return nil, err
	}
	departures := make([]Departure, len(ttssDepartures.Actual))
	for i, dep := range ttssDepartures.Actual {
		departures[i].PatternText = dep.PatternText
		departures[i].Direction = dep.Direction
		departures[i].PlannedTime = dep.PlannedTime
		departures[i].RelativeTime = dep.ActualRelativeTime
		departures[i].Predicted = dep.Status == "PREDICTED"
		// departures[i].Type = c.endpointType
	}
	return departures, nil
}
