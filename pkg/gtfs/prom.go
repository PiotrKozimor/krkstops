package gtfs

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	gtfsErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "gtfs_errors",
		Help: "Errors encountered in the gtfs package",
	}, []string{
		"err",
	})
)

func init() {
	prometheus.MustRegister(gtfsErrors)
}

func MetricHandler() http.Handler {
	return promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})
}
