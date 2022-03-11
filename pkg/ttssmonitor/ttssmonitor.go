package ttssmonitor

import (
	"errors"

	"github.com/PiotrKozimor/krkstops/pkg/ttss"
	"github.com/sirupsen/logrus"
)

// Errors below are sorted by severity - we care the most on the request that succeded,
// but returned no departures
const (
	OK               = iota
	OTHER_ERROR      // Please check logs for detailed error
	REQUEST_FAILED   // Request failed
	REQUEST_NON_200  // Request returned non 200 status code
	EMPTY_DEPARTURES // No departures received (with 200 status code)
	STOP_NOT_FOUND   // Stop was not found
)

func logIt(endpointId string, stopId uint) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"endpoint": endpointId,
		"stop":     stopId,
	})
}

// GetDeparturesErrorCode returns error that can be consumed by e.g. Prometheus.
func GetDeparturesErrorCode(e ttss.Endpointer, id uint) int {
	dep, err := e.GetDepartures(id)
	if err != nil {
		logIt(e.Id(), id).Error(err)
		if errors.Is(err, ttss.ErrStopNotFound) {
			return STOP_NOT_FOUND
		}
		if errors.Is(err, ttss.ErrRequestFailed) {
			return REQUEST_FAILED
		}
		if errors.Is(err, ttss.ErrStatusCode) {
			return REQUEST_NON_200
		}
		return OTHER_ERROR
	} else if len(dep) == 0 {
		logIt(e.Id(), id).Error("got no departures")
		return EMPTY_DEPARTURES
	}
	return OK
}
