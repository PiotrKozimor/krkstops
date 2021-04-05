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

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

func busHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/internetservice/services/passageInfo/stopPassages/stop":
		if r.URL.Query()["stop"][0] == "81" {
			w.WriteHeader(200)
			w.Write(mockBusDeparture81)
		} else if r.URL.Query()["stop"][0] == "610" {
			w.WriteHeader(200)
			w.Write(mockBusDeparture610)
		}
	case "/internetservice/geoserviceDispatcher/services/stopinfo/stops":
		w.WriteHeader(200)
		w.Write(mockStopsBus)
	}
}

func tramHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/internetservice/services/passageInfo/stopPassages/stop":
		if r.URL.Query()["stop"][0] == "81" {
			w.WriteHeader(400)
		} else if r.URL.Query()["stop"][0] == "610" {
			w.WriteHeader(200)
			w.Write(mockTramDeparture610)
		}
	case "/internetservice/geoserviceDispatcher/services/stopinfo/stops":
		w.WriteHeader(200)
		w.Write(mockStopsTram)
	}
}

func Ttss(ctx context.Context) {
	srvBus := http.Server{
		Addr:    "localhost:8080",
		Handler: HandlerFunc(busHandler),
	}
	srvTram := http.Server{
		Addr:    "localhost:8070",
		Handler: HandlerFunc(tramHandler),
	}
	go func() {
		err := srvBus.ListenAndServe()
		log.Print(err)
	}()
	go func() {
		err := srvTram.ListenAndServe()
		log.Print(err)
	}()
	<-ctx.Done()
	srvBus.Shutdown(ctx)
	srvTram.Shutdown(ctx)
	log.Printf("airly server shutdown")
}
