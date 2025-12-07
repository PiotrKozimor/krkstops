package gtfs

import (
	"log"
	"time"
)

var (
	Location *time.Location
)

func init() {
	var err error
	Location, err = time.LoadLocation("Europe/Warsaw")
	if err != nil {
		log.Fatal("load location: ", err)
	}
}
