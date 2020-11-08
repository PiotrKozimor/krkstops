package krkstops

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
	"github.com/RediSearch/redisearch-go/redisearch"
)

// TtssListStopsURL deccribes TTSS API endpoints used to fetch all stops. Separate endpoints exists for buses and trams.
var TtssListStopsURL = [...]string{
	BUS:  "http://91.223.13.70:80/internetservice/geoserviceDispatcher/services/stopinfo/stops",
	TRAM: "http://185.70.182.51:80/internetservice/geoserviceDispatcher/services/stopinfo/stops",
}

// TtssStops describes all stops returned by TTSS API
type TtssStops struct {
	Stops []pb.Stop
}

// StopsMap ShortName to Name
type StopsMap map[string]string

// GetAllStopsByURL fetches Stops from given endpoint.
func (app *App) GetAllStopsByURL(endp Endpoint) (StopsMap, error) {
	if int(endp) >= len(TtssListStopsURL) {
		return nil, fmt.Errorf("invalid endpoit provided: %d", endp)
	}
	url := TtssListStopsURL[endp]
	var stops = make(StopsMap, 500)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return stops, err
	}
	q := req.URL.Query()
	q.Add("left", "-648000000")
	q.Add("bottom", "-324000000")
	q.Add("right", "648000000")
	q.Add("top", "324000000")
	req.URL.RawQuery = q.Encode()
	resp, err := app.HTTPClient.Do(req)
	if err != nil {
		return stops, err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		if err != nil {
			return stops, err
		}
		return stops, errors.New(string(body))
	}
	var stopsWrapped TtssStops
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&stopsWrapped)
	if err != nil {
		return stops, err
	}
	for _, stop := range stopsWrapped.Stops {
		stops[stop.ShortName] = stop.Name
	}
	err = resp.Body.Close()
	return stops, err
}

// GetAllStops returns bus and tram stops.
func (app *App) GetAllStops() (StopsMap, error) {
	var stopsMap StopsMap = make(StopsMap, 1000)
	wg := sync.WaitGroup{}
	for endp := range TtssListStopsURL {
		wg.Add(1)
		go func(endp Endpoint) {
			defer wg.Done()
			stops, err := app.GetAllStopsByURL(endp)
			if err != nil {
				log.Println(err)
			}
			for shortName, name := range stops {
				stopsMap[shortName] = name
			}
		}(Endpoint(endp))
	}
	wg.Wait()
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
		pipe.Set(ShortName, name, -1)
		splitted := strings.Split(name, " ")
		for index, value := range splitted {
			if value == "(nż)" {
				splitted = append(splitted[:index], splitted[index+1:]...)
			}
		}
		for i := 0; i < len(splitted); i++ {
			swapped := append(splitted[i:], splitted[:i]...)
			term := strings.Join(swapped, " ")
			err := app.RedisAutocompleter.AddTerms(redisearch.Suggestion{Term: term, Score: 1, Payload: ShortName})
			if err != nil {
				return err
			}
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
