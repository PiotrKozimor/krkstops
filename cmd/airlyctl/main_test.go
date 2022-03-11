package main

import (
	"testing"

	"github.com/PiotrKozimor/krkstops/pkg/airly"
	"github.com/PiotrKozimor/krkstops/test"
)

func TestRootCmd(t *testing.T) {
	endpoint = airly.Endpoint("http://172.24.0.101:8072")
	test.Cmd(
		t,
		[]string{},
		rootCmdOutput,
		rootCmd,
	)
}

var rootCmdOutput = `CAQI  HUMIDITY[%]  TEMP [°C]  COLOR
12    46           12.3       6BC926`
