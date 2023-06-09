package ttss

import (
	"testing"

	"github.com/matryer/is"
)

func TestDepartures(t *testing.T) {
	is := is.New(t)
	cli := NewClient(Tram)
	deps, err := cli.GetDepartures(1360)
	is.NoErr(err)
	t.Log(deps)
}
