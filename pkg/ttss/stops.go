package ttss

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const stopsPath = "internetservice/geoserviceDispatcher/services/stopinfo/stops?left=-648000000&bottom=-324000000&right=648000000&top=324000000"

// TtssStops describes all stops returned by TTSS API
type ttssStops struct {
	Stops []stop
}

type stop struct {
	Name string
	Uid  uint32 `json:"shortName,string"`
}

// GetAllStops fetches Stops from given endpoint.
func (c Client) GetAllStops() ([]Stop, error) {
	resp, err := http.DefaultClient.Get(strings.Join([]string{c.host, stopsPath}, "/"))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("%w: %d", ErrStatusCode, resp.StatusCode)
		}
	}
	var stops ttssStops
	// io.Copy(os.Stdout, resp.Body)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&stops)
	parsedStops := make([]Stop, len(stops.Stops))
	for i, stop := range stops.Stops {
		parsedStops[i] = Stop{
			Name: stop.Name,
			Id:   uint(stop.Uid),
		}
	}
	return parsedStops, err
}

// // GetAllStops returns stops from multiple endpoints.
// // When one endpoint fails, valid stops are returned and error is send via chan.
// // When request is finished, error channel is closed.
// func GetAllStops(e []Endpointer) (chan []pb.Stop, chan error) {
// 	errC := make(chan error, len(e))
// 	stopsC := make(chan []pb.Stop, len(e))
// 	wg := sync.WaitGroup{}
// 	wg.Add(len(e))
// 	for _, endpoint := range e {
// 		go func(endp Endpointer) {
// 			stops, err := endp.GetAllStops()
// 			stopsC <- stops
// 			if err != nil {
// 				errC <- err
// 			}
// 			wg.Done()

// 		}(endpoint)
// 	}
// 	go func() {
// 		wg.Wait()
// 		close(errC)
// 		close(stopsC)
// 	}()
// 	return stopsC, errC
// }
