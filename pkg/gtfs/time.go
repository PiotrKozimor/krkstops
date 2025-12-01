package gtfs

import (
	"log"
	"time"
)

var (
	location *time.Location
)

func init() {
	var err error
	location, err = time.LoadLocation("Europe/Warsaw")
	if err != nil {
		log.Fatal("load location: ", err)
	}
}
