package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PiotrKozimor/krkstops/pkg/ttss"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var stopsPool = []uint{
	610,
	125,
}

var errorCodes = map[string]int{
	ttss.BusEndpoint.Id():  0,
	ttss.TramEndpoint.Id(): 0,
}

var (
	BUS  = ttss.BusEndpoint.Id()
	TRAM = ttss.TramEndpoint.Id()
)

func main() {
	l, err := logrus.ParseLevel(os.Getenv("LOGLEVEL"))
	if err != nil {
		log.Fatalln(err)
	}
	logrus.SetLevel(l)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	gauges := [...]prometheus.GaugeFunc{
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "departures_error_code",
			ConstLabels: prometheus.Labels{
				"endpoint": BUS,
			},
		}, func() float64 {
			errorCodeRet := errorCodes[BUS]
			logrus.WithFields(logrus.Fields{
				"endpoint": BUS,
				"code":     errorCodeRet,
			}).Info("Prometheus pull")
			errorCodes[BUS] = 0
			return float64(errorCodeRet)
		}),
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "departures_error_code",
			ConstLabels: prometheus.Labels{
				"endpoint": TRAM,
			},
		}, func() float64 {
			errorCodeRet := errorCodes[TRAM]
			logrus.WithFields(logrus.Fields{
				"endpoint": TRAM,
				"code":     errorCodeRet,
			}).Info("Prometheus pull")
			errorCodes[TRAM] = 0
			return float64(errorCodeRet)
		}),
	}
	r := prometheus.NewRegistry()
	for _, gauge := range gauges {
		r.MustRegister(gauge)
	}
	t := time.NewTicker(time.Minute * 2)
	go func() {
		for {
			<-t.C
			for _, endpoint := range ttss.KrkStopsEndpoints {
				for _, stop := range stopsPool {
					errCodeTmp := ttss.GetDeparturesErrorCode(endpoint, stop)
					// Any successfull request for Endpoint will erase other errors
					if errCodeTmp == 0 {
						errorCodes[endpoint.Id()] = 0
					} else if errCodeTmp > errorCodes[endpoint.Id()] {
						errorCodes[endpoint.Id()] = errCodeTmp
					}
				}
			}
		}
	}()
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	logrus.Warn("Ttssmonitord started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
