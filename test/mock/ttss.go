package mock

import (
	"context"
	_ "embed"
	"log"
	"net/http"
)

//go:embed data/bus1.json
var mockBusDeparture81 []byte

//go:embed data/bus2.json
var mockBusDeparture610 []byte

//go:embed data/tram2.json
var mockTramDeparture610 []byte

//go:embed data/stops_bus.json
var mockStopsBus []byte

//go:embed data/stops_tram.json
var mockStopsTram []byte

type handlerFunc func(http.ResponseWriter, *http.Request)

func (f handlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

func busHandler(w http.ResponseWriter, r *http.Request) {
	var resp []byte
	switch r.URL.Path {
	case "/internetservice/services/passageInfo/stopPassages/stop":
		if r.URL.Query()["stop"][0] == "81" {
			resp = mockBusDeparture81
		} else if r.URL.Query()["stop"][0] == "610" {
			resp = mockBusDeparture610
		}
	case "/internetservice/geoserviceDispatcher/services/stopinfo/stops":
		resp = mockStopsBus
	}
	writeResp(w, resp)
}

func tramHandler(w http.ResponseWriter, r *http.Request) {
	var resp []byte
	switch r.URL.Path {
	case "/internetservice/services/passageInfo/stopPassages/stop":
		if r.URL.Query()["stop"][0] == "81" {
			w.WriteHeader(400)
		} else if r.URL.Query()["stop"][0] == "610" {
			resp = mockTramDeparture610
		}
	case "/internetservice/geoserviceDispatcher/services/stopinfo/stops":
		resp = mockStopsTram
	}
	if resp != nil {
		writeResp(w, resp)
	}
}

func Ttss(ctx context.Context) {
	srvBus := http.Server{
		Addr:    "0.0.0.0:8071",
		Handler: handlerFunc(busHandler),
	}
	srvTram := http.Server{
		Addr:    "0.0.0.0:8070",
		Handler: handlerFunc(tramHandler),
	}
	go func() {
		err := srvBus.ListenAndServe()
		log.Print(err)
	}()
	err := srvTram.ListenAndServe()
	log.Print(err)
}
