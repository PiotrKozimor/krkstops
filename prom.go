package krkstops

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	registry     = prometheus.NewRegistry()
	serverErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "server_errors",
		Help: "Errors encountered in the server",
	}, []string{
		"err",
	})
	serverSuccesses = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "server_successes",
		Help: "Successful operations encountered in the server",
	}, []string{
		"name",
	})
	requests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests",
		Help: "HTTP requests completed",
	}, []string{
		"status_code",
		"path",
	})
)

func init() {
	registry.MustRegister(serverErrors, serverSuccesses, requests)
}

func MetricHandler() http.Handler {
	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
}
