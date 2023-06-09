package airly

import (
	"testing"

	"github.com/matryer/is"
)

var testEndpoint = Endpoint("http://172.24.0.101:8072")

func TestGetMeasurement(t *testing.T) {
	is := is.New(t)
	cli := NewClient(URL)
	m, err := cli.GetMeasurement(9895)
	is.NoErr(err)
	t.Logf("%+v", m)
}
