package ttssmonitor

import (
	"errors"
	"fmt"
	"testing"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/pkg/ttss"
	"github.com/matryer/is"
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
	is := is.New(t)
	c := make(chan error, 4)
	te := TestEndpoint{
		c: c,
	}
	c <- fmt.Errorf("%w: %v", ttss.ErrRequestFailed, "err")
	c <- fmt.Errorf("%w: %d", ttss.ErrStatusCode, 400)
	c <- errors.New("bar")
	c <- nil
	is.Equal(REQUEST_FAILED, GetDeparturesErrorCode(te, 0))
	is.Equal(REQUEST_NON_200, GetDeparturesErrorCode(te, 0))
	is.Equal(OTHER_ERROR, GetDeparturesErrorCode(te, 0))
	is.Equal(EMPTY_DEPARTURES, GetDeparturesErrorCode(te, 0))
}
