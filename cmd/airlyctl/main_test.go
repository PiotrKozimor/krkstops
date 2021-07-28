package main

import (
	"testing"

	"github.com/PiotrKozimor/krkstops/test"
)

func TestRootCmd(t *testing.T) {
	test.Cmd(
		t,
		[]string{},
		rootCmdOutput,
		rootCmd,
	)
}

var rootCmdOutput = `CAQI  HUMIDITY[%]  TEMP [°C]  COLOR
12    46           12.3       6BC926`
