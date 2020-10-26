package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PiotrKozimor/krk-stops-backend-golang/krkstops"
	pb "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
	"github.com/PiotrKozimor/krk-stops-backend-golang/ttssmonitor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var stopsPool = []pb.Stop{
	{
		ShortName: "610",
	},
	{
		ShortName: "125",
	},
}

var errorCodes = []int{
	krkstops.BUS:  0,
	krkstops.TRAM: 0,
}

func main() {
	l, err := logrus.ParseLevel(os.Getenv("LOGLEVEL"))
	if err != nil {
		log.Fatalln(err)
	}
	logrus.SetLevel(l)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	for endpoint := range errorCodes {
		var gauge = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name: fmt.Sprintf("departures_error_code_%d", endpoint),
			Help: fmt.Sprintf("Current error code of TTSS API for %d", endpoint),
		}, func() float64 {
			errorCodeRet := errorCodes[endpoint]
			logrus.WithFields(logrus.Fields{
				"endpoint": endpoint,
				"code":     errorCodeRet,
			}).Info("Prometheus pull")
			errorCodes[endpoint] = 0
			return float64(errorCodeRet)
		})
		prometheus.MustRegister(gauge)
	}
	t := time.NewTicker(time.Minute * 2)
	go func() {
		for {
			<-t.C
			for endpoint := range krkstops.TtssStopDearturesURLs {
				for _, stop := range stopsPool {
					errCodeTmp := ttssmonitor.GetStopDeparturesErrorCode(krkstops.Endpoint(endpoint), &stop)
					// Any successfull request for Endpoint will erase other errors
					if errCodeTmp == 0 {
						errorCodes[endpoint] = 0
					} else if errCodeTmp > errorCodes[endpoint] {
						errorCodes[endpoint] = errCodeTmp
					}
				}
			}
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	logrus.Warn("Ttssmonitord started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
