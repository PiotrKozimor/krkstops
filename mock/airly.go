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
	switch r.URL.Path {
	case "/v2/installations/8077":
		w.WriteHeader(200)
		w.Write(mockInstallation)
	case "/v2/installations/nearest":
		w.WriteHeader(200)
		w.Write(mockNearest)
	case "/v2/measurements/installation":
		w.WriteHeader(200)
		w.Write(mockMeasurement)
	}
}

func Airly(ctx context.Context) {
	srv := http.Server{
		Addr:    "0.0.0.0:8072",
		Handler: HandlerFunc(airlyHandler),
	}
	go func() {
		err := srv.ListenAndServe()
		log.Print(err)
	}()
	<-ctx.Done()
	srv.Shutdown(ctx)
}
