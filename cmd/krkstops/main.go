package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	_ "time/tzdata"

	"github.com/PiotrKozimor/krkstops"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	server, err := krkstops.NewServer()
	handle(err)

	go func() {
		log.Printf("http metrics server listening on :8083")
		handle(http.ListenAndServe(":8083", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})))
	}()

	log.Printf("http rpc server listening on :8082")
	handle(http.ListenAndServe(":8082", server))
}
