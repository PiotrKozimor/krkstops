package ttss

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/PiotrKozimor/krkstops/pb"
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

// Stops ShortName to Name
type Stops map[string]string

// GetAllStopsByURL fetches Stops from given endpoint.
func GetAllStopsByURL(endp Endpoint) (Stops, error) {
	if int(endp) >= len(TtssListStopsURL) {
		return nil, fmt.Errorf("invalid endpoit provided: %d", endp)
	}
	url := TtssListStopsURL[endp]
	var stops = make(Stops, 500)
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
	resp, err := http.DefaultClient.Do(req)
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
func GetAllStops() (Stops, error) {
	var stopsMap Stops = make(Stops, 1000)
	wg := sync.WaitGroup{}
	for endp := range TtssListStopsURL {
		wg.Add(1)
		go func(endp Endpoint) {
			defer wg.Done()
			stops, err := GetAllStopsByURL(endp)
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
