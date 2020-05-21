package krkstops

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v7"
)

const (
	tramTtssStopDearturesURL = "http://185.70.182.51/internetservice/services/passageInfo/stopPassages/stop"
	busTtssStopDearturesURL  = "http://91.223.13.70/internetservice/services/passageInfo/stopPassages/stop"
)

// Departure describe one departure - tram on bus from given stop.
type departure struct {
	PlannedTime        string
	Status             string
	ActualRelativeTime int32
	Direction          string
	PatternText        string
}

// StopDepartures includes all departures from given Stop.
type stopDepartures struct {
	Actual []departure
}

// App stores clients
type App struct {
	HTTPClient         *http.Client
	RedisClient        *redis.Client
	RedisAutocompleter *redisearch.Autocompleter
}

// GetStopDepartures returns StopDepartures from given Stop.
func (app *App) getStopDeparturesByURL(url string, stop *pb.Stop, c chan []pb.Departure) error {
	var departures []pb.Departure
	var departure pb.Departure
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c <- departures
		return err
	}
	q := req.URL.Query()
	q.Add("stop", fmt.Sprint(stop.ShortName))
	q.Add("mode", "departure")
	q.Add("language", "pl")
	req.URL.RawQuery = q.Encode()
	resp, err := app.HTTPClient.Do(req)
	if err != nil {
		c <- departures
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c <- departures
		return err
	}
	if resp.StatusCode != 200 {
		c <- departures
		return errors.New(string(body))
	}
	var stopDepartures stopDepartures
	err = json.Unmarshal(body, &stopDepartures)
	if err != nil {
		c <- departures
		return err
	}
	for _, dep := range stopDepartures.Actual {
		departure.PatternText = dep.PatternText
		departure.Direction = dep.Direction
		departure.PlannedTime = dep.PlannedTime
		if dep.Status == "PREDICTED" {
			departure.RelativeTime = dep.ActualRelativeTime
		} else {
			departure.RelativeTime = 0
		}
		departures = append(departures, departure)
	}
	c <- departures
	return nil
}

// GetStopDepartures returns bus and tram StopDepartures from given Stop.
func (app *App) GetStopDepartures(stop *pb.Stop) ([]pb.Departure, error) {
	c := make(chan []pb.Departure)
	go app.getStopDeparturesByURL(busTtssStopDearturesURL, stop, c)
	go app.getStopDeparturesByURL(tramTtssStopDearturesURL, stop, c)
	deps1, deps2 := <-c, <-c
	deps := append(deps1, deps2...)
	sort.Slice(deps, func(i, j int) bool {
		return deps[i].PlannedTime < deps[j].PlannedTime
	})
	return deps, nil
}
