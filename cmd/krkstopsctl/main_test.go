package main

import (
	"testing"

	"github.com/PiotrKozimor/krkstops"
	"github.com/PiotrKozimor/krkstops/test"
	"github.com/matryer/is"
)

func TestAirlyCmd(t *testing.T) {
	test.Cmd(
		t,
		[]string{
			"-e",
			"localhost:8080",
			"airly",
		},
		airlyCmdOutput,
		rootCmd,
	)
}

func TestDepsCmd(t *testing.T) {
	test.Cmd(
		t,
		[]string{
			"-e",
			"localhost:8080",
			"deps",
		},
		depsCmdOutput,
		rootCmd,
	)
}

func TestStopsCmd(t *testing.T) {
	is := is.New(t)
	cache, err := krkstops.NewScore("localhost:6379", krkstops.SUG)
	is.NoErr(err)
	err = cache.Update()
	is.NoErr(err)
	test.Cmd(
		t,
		[]string{
			"-e",
			"localhost:8080",
			"stops",
			"mat",
		},
		stopsCmdOutput,
		rootCmd,
	)
}

var airlyCmdOutput = `CAQI  HUMIDITY[%]  TEMP [°C]  COLOR
12    46           12.3       6BC926`

var depsCmdOutput = `
179  Os. Kurdwanów            36         14:48
304  Wieliczka Miasto         336        14:55
173  Nowy Bieżanów Południe   396        14:55
164  Górka Narodowa           456        14:57
503  Górka Narodowa Wschód    456        14:55
304  Dworzec Główny Zachód    636        15:00
169  Zajezdnia Wola Duchacka  756        14:59
179  Dworzec Główny Zachód    996        15:06
144  Azory                    1056       15:06
503  Nowy Bieżanów Południe   1056       15:07
164  Piaski Nowe              1356       15:11
173  Azory                    1476       15:14
304  Wieliczka Miasto         1536       15:15
169  Górka Narodowa           1536       15:15
10   Pleszów                  156        14:52
10   Łagiewniki               456        14:56
11   Czerwone Maki P+R        516        14:58
11   Mały Płaszów P+R         696        15:00
8    Borek Fałęcki            816        15:03
8    Bronowice Małe           1116       15:08
10   Pleszów                  1356       15:12`

var stopsCmdOutput = `NO  ID   NAME
0   610  Rondo Matecznego`
