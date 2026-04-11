package krkstops

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	serverErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "server_errors",
		Help: "Errors encountered in the server (logged)",
	}, []string{
		"err",
	})
	serverSuccesses = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "server_successes",
		Help: "Successful operations encountered in the server (not logged)",
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
	prometheus.DefaultRegisterer.MustRegister(serverErrors, serverSuccesses, requests)
}
