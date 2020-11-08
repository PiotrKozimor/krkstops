package krkstops

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v7"
)

type Endpoint int

const (
	BUS Endpoint = iota
	TRAM
	INVALID
)

// TtssStopDearturesURLs deccribes TTSS API endpoints used to query departures. Separate endpoints exists for buses and trams.
var TtssStopDearturesURLs = [...]string{
	BUS:  "http://91.223.13.70/internetservice/services/passageInfo/stopPassages/stop",
	TRAM: "http://185.70.182.51/internetservice/services/passageInfo/stopPassages/stop",
}

// Departure describe one Departure - tram or bus from given stop.
type Departure struct {
	PlannedTime        string
	Status             string
	ActualRelativeTime int32
	Direction          string
	PatternText        string
}

// StopDepartures includes all departures from given Stop.
type StopDepartures struct {
	Actual []Departure
}

// App stores clients
type App struct {
	HTTPClient         *http.Client
	RedisClient        *redis.Client
	RedisAutocompleter *redisearch.Autocompleter
}

// GetStopDeparturesByURL fetches Departures from given Stop from given endpoint.
func (app *App) GetStopDeparturesByURL(endp Endpoint, stop *pb.Stop) ([]pb.Departure, error) {
	if int(endp) >= len(TtssStopDearturesURLs) {
		return nil, fmt.Errorf("invalid endpoit provided: %d", endp)
	}
	departures := make([]pb.Departure, 0, 20)
	var departure pb.Departure
	req, err := http.NewRequest("GET", TtssStopDearturesURLs[endp], nil)
	if err != nil {
		return departures, err
	}
	q := req.URL.Query()
	q.Add("stop", fmt.Sprint(stop.ShortName))
	q.Add("mode", "departure")
	q.Add("language", "pl")
	req.URL.RawQuery = q.Encode()
	resp, err := app.HTTPClient.Do(req)
	if err != nil {
		return departures, err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		if err != nil {
			return departures, err
		}
		return departures, errors.New(string(body))
	}
	var stopDepartures StopDepartures
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&stopDepartures)
	if err != nil {
		return departures, err
	}
	for _, dep := range stopDepartures.Actual {
		departure.PatternText = dep.PatternText
		departure.Direction = dep.Direction
		departure.PlannedTime = dep.PlannedTime
		departure.RelativeTime = dep.ActualRelativeTime
		departure.Predicted = dep.Status == "PREDICTED"
		departures = append(departures, departure)
	}
	err = resp.Body.Close()
	return departures, err
}

// GetStopDepartures returns bus and tram StopDepartures from given Stop.
func (app *App) GetStopDepartures(stop *pb.Stop) ([]pb.Departure, error) {
	departures := make([]pb.Departure, 0, 40)
	c := make(chan []pb.Departure)
	for endp := range TtssStopDearturesURLs {
		go func(endp Endpoint) {
			deps, err := app.GetStopDeparturesByURL(endp, stop)
			if err != nil {
				log.Println(err)
			}
			c <- deps
		}(Endpoint(endp))
	}
	for range TtssStopDearturesURLs {
		depsTmp := <-c
		departures = append(departures, depsTmp...)
	}
	sort.Slice(departures, func(i, j int) bool {
		return departures[i].RelativeTime < departures[j].RelativeTime
	})
	return departures, nil
}
