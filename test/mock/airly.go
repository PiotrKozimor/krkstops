package mock

import (
	"context"
	_ "embed"
	"log"
	"net/http"
)

//go:embed data/installation.json
var mockInstallation []byte

//go:embed data/measurement.json
var mockMeasurement []byte

//go:embed data/nearest.json
var mockNearest []byte

// var ctx, cancel = context.WithTimeout(context.Background(), time.Second)

func airlyHandler(w http.ResponseWriter, r *http.Request) {
	var resp []byte
	switch r.URL.Path {
	case "/v2/installations/8077":
		resp = mockInstallation
	case "/v2/installations/nearest":
		resp = mockNearest
	case "/v2/measurements/installation":
		resp = mockMeasurement
	}
	writeResp(w, resp)
}

func writeResp(rw http.ResponseWriter, resp []byte) {
	rw.WriteHeader(200)
	_, err := rw.Write(resp)
	if err != nil {
		log.Println("airly handler: ", err)
	}
}

func Airly(ctx context.Context) {
	srv := http.Server{
		Addr:    "0.0.0.0:8072",
		Handler: handlerFunc(airlyHandler),
	}
	err := srv.ListenAndServe()
	log.Print(err)
}
