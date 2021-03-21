package ttss

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PiotrKozimor/krkstops/pb"
)

type Endpoint string

type ErrStatusCode struct {
	code int
}

type ErrRequestFailed struct {
	err error
}

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

type Endpointer interface {
	GetDepartures(uint) ([]pb.Departure, error)
	GetAllStops() ([]pb.Stop, error)
	Id() string
}

const (
	departuresPath = "services/passageInfo/stopPassages/stop?stop=%d&mode=departure&language=pl"
	BusEndpoint    = Endpoint("http://91.223.13.70/internetservice")
	TramEndpoint   = Endpoint("http://185.70.182.51/internetservice")
)

var endpointsIds = map[Endpoint]string{
	BusEndpoint:  "bus",
	TramEndpoint: "tram",
}

var KrkStopsEndpoints = []Endpointer{
	BusEndpoint,
	TramEndpoint,
}

func (e ErrStatusCode) Error() string {
	return fmt.Sprintf("status code: %d", e.code)
}

func (e ErrRequestFailed) Error() string {
	return fmt.Sprintf("request failed: %v", e.err)
}

func (e Endpoint) Id() string {
	return endpointsIds[e]
}

// GetDepartures from Endpoint for stop with given shortName.
func (e Endpoint) GetDepartures(shortName uint) ([]pb.Departure, error) {
	resp, err := http.DefaultClient.Get(fmt.Sprintf(strings.Join([]string{string(e), departuresPath}, "/"), shortName))
	if err != nil {
		return nil, ErrRequestFailed{err: err}
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, ErrStatusCode{code: resp.StatusCode}
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
			if err != nil {
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
