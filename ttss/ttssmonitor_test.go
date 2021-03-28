package ttss

import (
	"errors"
	"testing"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/stretchr/testify/assert"
)

type TestEndpoint struct {
	c chan error
}

func (e TestEndpoint) GetDepartures(uint) ([]pb.Departure, error) {
	return nil, <-e.c
}

func (e TestEndpoint) GetAllStops() (a []pb.Stop, b error) { return }
func (e TestEndpoint) Id() (a string)                      { return }

func TestGetDeparturesErrorCode(t *testing.T) {
	c := make(chan error, 4)
	te := TestEndpoint{
		c: c,
	}
	c <- ErrRequestFailed{err: errors.New("test")}
	c <- ErrStatusCode{code: 401}
	c <- errors.New("bar")
	c <- nil
	assert.Equal(t, REQUEST_FAILED, GetDeparturesErrorCode(te, 0))
	assert.Equal(t, REQUEST_NON_200, GetDeparturesErrorCode(te, 0))
	assert.Equal(t, OTHER_ERROR, GetDeparturesErrorCode(te, 0))
	assert.Equal(t, EMPTY_DEPARTURES, GetDeparturesErrorCode(te, 0))
}
