package search

import (
	"testing"

	"github.com/matryer/is"
)

func TestSearch(t *testing.T) {
	is := is.New(t)

	s, err := New()
	is.NoErr(err)

	// s.t.Dump()
	for _, str := range []string{
		"rón",
		"zav",
		"zac",
		"la",
		"łe",
		"le",
		"os",
		"ks",
		"m1",
	} {
		res := s.Search(str, 1e6)
		t.Logf("%s", str)
		for _, stop := range res {
			t.Logf("\t%d\t%s", stop.Score, stop.Name)
		}
	}
}

func TestExact(t *testing.T) {
	is := is.New(t)

	s, err := New()
	is.NoErr(err)

	for _, str := range []string{
		"zac",
		"no",
	} {
		res := s.SearchExact(str, 1e6)
		t.Logf("%s", str)
		for _, stop := range res {
			t.Logf("\t%d\t%s", stop.Score, stop.Name)
		}
	}
}

func TestWithinDistance(t *testing.T) {
	is := is.New(t)

	s, err := New()
	is.NoErr(err)

	for _, str := range []string{
		"rón",
		"zav",
		"la",
	} {
		res := s.searchWithinDistance(str, 1e6)
		t.Logf("%s", str)
		for _, stop := range res {
			t.Logf("\t%d\t%s", stop.Score, stop.Name)
		}
	}
}
