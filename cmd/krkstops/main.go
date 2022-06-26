package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strings"

	"github.com/PiotrKozimor/krkstops"
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/pkg/airly"
	"github.com/PiotrKozimor/krkstops/pkg/ttss"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var stop = make(chan bool)

type Server interface {
	Serve(net.Listener) error
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func serveGracefully(s Server, l net.Listener) {
	err := s.Serve(l)
	if err != nil {
		log.Print(err)
		stop <- true
	}
}

type Config struct {
	RedisURI          string
	TlsCert           string
	TlsKey            string
	OverrideEndpoints string
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	handle(err)
	var conf Config
	err = envconfig.Process("", &conf)
	log.Printf("CONFIG: %+v\n", conf)
	handle(err)
	grpcServer := grpc.NewServer()
	var grpcServerTls *grpc.Server
	server := krkstops.NewServer(conf.RedisURI)
	server.Airly, server.Ttss = Endpoints(conf.OverrideEndpoints)
	pb.RegisterKrkStopsServer(grpcServer, server)
	tlsOn := conf.TlsCert != "" && conf.TlsKey != ""
	if tlsOn {
		cert, err := tls.LoadX509KeyPair(conf.TlsCert, conf.TlsKey)
		handle(err)
		grpcServerTls = grpc.NewServer(
			grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		)
		pb.RegisterKrkStopsServer(grpcServerTls, server)
	}
	if tlsOn {
		lisTls, err := net.Listen("tcp", ":9090")
		handle(err)
		go serveGracefully(grpcServerTls, lisTls)
		defer grpcServer.GracefulStop()
	}
	go serveGracefully(grpcServer, lis)
	defer grpcServer.GracefulStop()
	go func() {
		HandlePrivacy()
		var err error
		if tlsOn {
			err = http.ListenAndServeTLS(":443", conf.TlsCert, conf.TlsKey, nil)
		} else {
			err = http.ListenAndServe(":8090", nil)
		}
		if err != nil {
			log.Print(err)
			stop <- true
		}
	}()
	<-stop
}

// If OVERRIDEENDPOINTS environmental variable is set with structure
// <airlyEndpoint>,<busEndpoint>,<tramEndpoint> the endpoints are overriden.
// Otherwise, the defaults are used.
func Endpoints(override string) (airly.Endpoint, []ttss.Endpointer) {
	if override != "" {
		endpoints := strings.Split(override, ",")
		return airly.Endpoint(endpoints[0]),
			[]ttss.Endpointer{
				ttss.Endpoint{
					URL:  endpoints[1],
					Type: pb.Endpoint_BUS,
				},
				ttss.Endpoint{
					URL:  endpoints[2],
					Type: pb.Endpoint_TRAM,
				},
			}
	}
	return airly.DefaultEndpoint, ttss.KrkStopsEndpoints
}
