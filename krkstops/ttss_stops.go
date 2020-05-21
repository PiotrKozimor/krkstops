package krkstops

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
	"github.com/RediSearch/redisearch-go/redisearch"
)

const (
	tramTTssListStopsURL = "http://185.70.182.51:80/internetservice/geoserviceDispatcher/services/stopinfo/stops"
	busTTssListStopsURL  = "http://91.223.13.70:80/internetservice/geoserviceDispatcher/services/stopinfo/stops"
)

// TtssStops describes all stops returned by TTSS API
type TtssStops struct {
	Stops []pb.Stop
}

// StopsMap ShortName to Name
type StopsMap map[string]string

func (app *App) getAllStopsByURL(url string, c chan []pb.Stop) error {
	var stps []pb.Stop
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c <- stps
		return err
	}
	q := req.URL.Query()
	q.Add("left", "-648000000")
	q.Add("bottom", "-324000000")
	q.Add("right", "648000000")
	q.Add("top", "324000000")
	// req.URL.RawQuery = "left=-648000000&bottom=-324000000&right=648000000&top=324000000"
	req.URL.RawQuery = q.Encode()
	resp, err := app.HTTPClient.Do(req)
	if err != nil {
		c <- stps
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	// log.Printf("%v", string(body))
	if err != nil {
		c <- stps
		return err
	}
	if resp.StatusCode != 200 {
		c <- stps
		return errors.New(string(body))
	}
	var stopsWrapped TtssStops
	err = json.Unmarshal(body, &stopsWrapped)
	if err != nil {
		c <- stps
		return err
	}
	c <- stopsWrapped.Stops
	return nil
}

// GetAllStops returns bus and tram stops.
func (app *App) GetAllStops() (StopsMap, error) {
	c := make(chan []pb.Stop)
	go app.getAllStopsByURL(busTTssListStopsURL, c)
	go app.getAllStopsByURL(tramTTssListStopsURL, c)
	stops1, stops2 := <-c, <-c
	allStops := append(stops1, stops2...)
	// log.Printf("%v", allStops)
	var stopsMap StopsMap = make(StopsMap)
	for _, stop := range allStops {
		stopsMap[stop.ShortName] = stop.Name
	}
	return stopsMap, nil
}

// CompareStops and return new and old. New stops are stored in 'stops.tmp' set.
func (app *App) CompareStops(stops StopsMap) (newStops StopsMap, oldStops StopsMap, err error) {
	newStops = make(StopsMap)
	oldStops = make(StopsMap)
	pipe := app.RedisClient.TxPipeline()
	for ShortName := range stops {
		pipe.SAdd("stops.tmp", ShortName)
	}
	_, err = pipe.Exec()
	if err != nil {
		return nil, nil, err
	}
	newStopsShortNames := app.RedisClient.SDiff("stops.tmp", "stops").Val()
	oldStopsShortNames := app.RedisClient.SDiff("stops", "stops.tmp").Val()
	for _, ShortName := range newStopsShortNames {
		newStops[ShortName] = stops[ShortName]
	}
	for _, ShortName := range oldStopsShortNames {
		oldStops[ShortName] = stops[ShortName]
	}
	return newStops, oldStops, nil
}

// UpdateSuggestionsAndRedis in Redisearch engine
func (app *App) UpdateSuggestionsAndRedis(newStops StopsMap, oldStops StopsMap) error {
	pipe := app.RedisClient.TxPipeline()
	for ShortName, name := range newStops {
		pipe.SAdd("stops.new", ShortName)
		err := app.RedisAutocompleter.AddTerms(redisearch.Suggestion{Term: name, Score: 1, Payload: ShortName})
		if err != nil {
			return err
		}
	}
	_, err := pipe.Exec()
	if err != nil {
		return err
	}
	for _, name := range oldStops {
		err := app.RedisAutocompleter.DeleteTerms(redisearch.Suggestion{Term: name})
		if err != nil {
			return err
		}
	}
	_, err = app.RedisClient.Rename("stops.tmp", "stops").Result()
	if err != nil {
		return err
	}
	return nil
}
