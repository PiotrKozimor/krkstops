package ttss

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PiotrKozimor/krkstops/pb"
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

func (e Endpoint) Id() string {
	return endpointsIds[e]
}

// GetDepartures from Endpoint for stop with given shortName.
func (e Endpoint) GetDepartures(id uint) ([]pb.Departure, error) {
	resp, err := http.DefaultClient.Get(fmt.Sprintf(strings.Join([]string{e.URL, departuresPath}, "/"), id))
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
	departures := make([]pb.Departure, len(ttssDepartures.Actual))
	for i, dep := range ttssDepartures.Actual {
		departures[i].PatternText = dep.PatternText
		departures[i].Direction = dep.Direction
		departures[i].PlannedTime = dep.PlannedTime
		departures[i].RelativeTime = dep.ActualRelativeTime
		departures[i].Predicted = dep.Status == "PREDICTED"
		departures[i].Type = e.Type
	}
	return departures, nil
}

// GetDepartures for range of endpoints. When one endpoint fails, valid departures are returned and error is send via chan.
// When request is finished, error channel is closed.
func GetDepartures(e []Endpointer, shortName uint) (chan []pb.Departure, chan error) {
	errC := make(chan error, len(e))
	depC := make(chan []pb.Departure, len(e))
	wg := sync.WaitGroup{}
	wg.Add(len(e))
	for _, endpoint := range e {
		go func(endp Endpointer) {
			departures, err := endp.GetDepartures(shortName)
			depC <- departures
			if !errors.Is(err, ErrStopNotFound) {
				errC <- err
			}
			wg.Done()

		}(endpoint)
	}
	go func() {
		wg.Wait()
		close(errC)
		close(depC)
	}()
	return depC, errC
}
