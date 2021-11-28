package main

import (
	"github.com/PiotrKozimor/krkstops/cmd/gtfs/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	cmd.Execute()
}
