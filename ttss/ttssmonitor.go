package ttss

import (
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
)

func log(endpointId string, stopId uint) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"endpoint": endpointId,
		"stop":     stopId,
	})
}

// GetDeparturesErrorCode returns error that can be consumed by e.g. Prometheus.
func GetDeparturesErrorCode(e Endpointer, id uint) (errorCode int) {
	dep, err := e.GetDepartures(id)
	if err != nil {
		switch err.(type) {
		case ErrRequestFailed:
			errorCode = REQUEST_FAILED
		case ErrStatusCode:
			errorCode = REQUEST_NON_200
		default:
			errorCode = OTHER_ERROR
		}
		log(e.Id(), id).Error(err)
	} else if len(dep) == 0 {
		errorCode = EMPTY_DEPARTURES
		log(e.Id(), id).Error("got no departures")
	}
	return
}
