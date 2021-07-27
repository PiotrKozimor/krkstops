package main

import (
	"testing"

	"github.com/PiotrKozimor/krkstops/test"
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
	test.Cmd(
		t,
		[]string{
			"-e",
			"localhost:8080",
			"stops",
		},
		depsCmdOutput,
		rootCmd,
	)
}

var airlyCmdOutput = `CAQI  HUMIDITY[%]  TEMP [°C]  COLOR
12    46           12.3       6BC926`

var depsCmdOutput = `relativeTime:156 plannedTime:"14:52" direction:"Pleszów" patternText:"10" predicted:true type:TRAM
relativeTime:456 plannedTime:"14:56" direction:"Łagiewniki" patternText:"10" predicted:true type:TRAM
relativeTime:516 plannedTime:"14:58" direction:"Czerwone Maki P+R" patternText:"11" predicted:true type:TRAM
relativeTime:696 plannedTime:"15:00" direction:"Mały Płaszów P+R" patternText:"11" predicted:true type:TRAM
relativeTime:816 plannedTime:"15:03" direction:"Borek Fałęcki" patternText:"8" type:TRAM
relativeTime:1116 plannedTime:"15:08" direction:"Bronowice Małe" patternText:"8" predicted:true type:TRAM
relativeTime:1356 plannedTime:"15:12" direction:"Pleszów" patternText:"10" predicted:true type:TRAM
relativeTime:36 plannedTime:"14:48" direction:"Os. Kurdwanów" patternText:"179" predicted:true
relativeTime:336 plannedTime:"14:55" direction:"Wieliczka Miasto" patternText:"304"
relativeTime:396 plannedTime:"14:55" direction:"Nowy Bieżanów Południe" patternText:"173" predicted:true
relativeTime:456 plannedTime:"14:57" direction:"Górka Narodowa" patternText:"164" predicted:true
relativeTime:456 plannedTime:"14:55" direction:"Górka Narodowa Wschód" patternText:"503" predicted:true
relativeTime:636 plannedTime:"15:00" direction:"Dworzec Główny Zachód" patternText:"304"
relativeTime:756 plannedTime:"14:59" direction:"Zajezdnia Wola Duchacka" patternText:"169" predicted:true
relativeTime:996 plannedTime:"15:06" direction:"Dworzec Główny Zachód" patternText:"179" predicted:true
relativeTime:1056 plannedTime:"15:06" direction:"Azory" patternText:"144" predicted:true
relativeTime:1056 plannedTime:"15:07" direction:"Nowy Bieżanów Południe" patternText:"503" predicted:true
relativeTime:1356 plannedTime:"15:11" direction:"Piaski Nowe" patternText:"164" predicted:true
relativeTime:1476 plannedTime:"15:14" direction:"Azory" patternText:"173" predicted:true
relativeTime:1536 plannedTime:"15:15" direction:"Wieliczka Miasto" patternText:"304"
relativeTime:1536 plannedTime:"15:15" direction:"Górka Narodowa" patternText:"169" predicted:true`
