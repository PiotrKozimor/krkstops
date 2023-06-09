package airly

import (
	"testing"

	"github.com/matryer/is"
)

func TestNearestInstallation(t *testing.T) {
	is := is.New(t)
	cli := NewClient(URL)
	i, err := cli.NearestInstallation(50.049880, 19.926861)
	is.NoErr(err)
	t.Logf("%+v", i)
}

func TestGetInstallation(t *testing.T) {
	is := is.New(t)
	cli := NewClient(URL)
	i, err := cli.GetInstallation(8077)
	is.NoErr(err)
	t.Logf("%+v", i)
}
