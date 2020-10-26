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

// TtssStopDearturesURLs deccribes TTSS API endpoints used to query departures. Separate endpoints exists for buses and trams.
var TtssStopDearturesURLs = map[string]string{
	"tram": "http://185.70.182.51/internetservice/services/passageInfo/stopPassages/stop",
	"bus":  "http://91.223.13.70/internetservice/services/passageInfo/stopPassages/stop",
}

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

// GetStopDeparturesByURL fetches Departures from given Stop from given endpoint.
func (app *App) GetStopDeparturesByURL(url string, stop *pb.Stop) ([]pb.Departure, error) {
	departures := make([]pb.Departure, 0, 20)
	var departure pb.Departure
	req, err := http.NewRequest("GET", url, nil)
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
	var stopDepartures stopDepartures
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
	return departures, nil
}

// GetStopDepartures returns bus and tram StopDepartures from given Stop.
func (app *App) GetStopDepartures(stop *pb.Stop) ([]pb.Departure, error) {
	departures := make([]pb.Departure, 0, 40)
	c := make(chan []pb.Departure)
	for _, url := range TtssStopDearturesURLs {
		go func(url string) {
			deps, err := app.GetStopDeparturesByURL(url, stop)
			if err != nil {
				log.Println(err)
			}
			c <- deps
		}(url)
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
