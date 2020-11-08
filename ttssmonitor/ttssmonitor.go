package ttssmonitor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/PiotrKozimor/krk-stops-backend-golang/krkstops"
	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
	log "github.com/sirupsen/logrus"
)

// Errors below are sorted by severity - we care the most on the request that succeded, but returned no departures
const (
	OK               = iota
	OTHER_ERROR      // Please check logs for detailed error
	REQUEST_FAILED   // Request failed
	REQUEST_NON_200  // Request returned non 200 status code
	EMPTY_DEPARTURES // No departures received (with 200 status code)
)

func logWithEndpointAndStop(endp krkstops.Endpoint, stop *pb.Stop) *log.Entry {
	return log.WithFields(log.Fields{
		"endpoint": endp,
		"stop":     stop.ShortName,
	})
}

// GetStopDeparturesErrorCodeByURL returns error that can be consumed by e.g. Prometheus.
func GetStopDeparturesErrorCode(endp krkstops.Endpoint, stop *pb.Stop) (errorCode int) {
	if int(endp) >= len(krkstops.TtssListStopsURL) {
		logWithEndpointAndStop(endp, stop).Errorln("invalid endpoit provided: %d", endp)
		return OTHER_ERROR
	}
	url := krkstops.TtssStopDearturesURLs[endp]
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logWithEndpointAndStop(endp, stop).Errorln(err)
		return OTHER_ERROR
	}
	q := req.URL.Query()
	q.Add("stop", fmt.Sprint(stop.ShortName))
	q.Add("mode", "departure")
	q.Add("language", "pl")
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logWithEndpointAndStop(endp, stop).Warn(err)
		return REQUEST_FAILED
	}
	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logWithEndpointAndStop(endp, stop).Error(err)
		}
		err = resp.Body.Close()
		if err != nil {
			logWithEndpointAndStop(endp, stop).Error(err)
		}
		logWithEndpointAndStop(endp, stop).Warnf("StatusCode: %d, Body: %s", resp.StatusCode, body)
		return REQUEST_NON_200
	}
	var stopDepartures krkstops.StopDepartures
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&stopDepartures)
	if err != nil {
		logWithEndpointAndStop(endp, stop).Error(err)
		return OTHER_ERROR
	}
	err = resp.Body.Close()
	if err != nil {
		logWithEndpointAndStop(endp, stop).Error(err)
		return OTHER_ERROR
	}
	if len(stopDepartures.Actual) == 0 {
		logWithEndpointAndStop(endp, stop).Warn(stopDepartures.Actual)
		return EMPTY_DEPARTURES
	}
	logWithEndpointAndStop(endp, stop).Infoln("Requests succeeded.")
	return OK
}
