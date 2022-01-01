package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"

	"github.com/PiotrKozimor/krkstops"
	"github.com/PiotrKozimor/krkstops/pb"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	RedisURI string
	TlsCert  string
	TlsKey   string
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	handle(err)
	var conf Config
	err = envconfig.Process("", &conf)
	log.Printf("CONFIG: %+v\n", conf)
	handle(err)
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	var grpcServerTls *grpc.Server
	server, err := krkstops.NewServer(conf.RedisURI)
	handle(err)
	pb.RegisterKrkStopsServer(grpcServer, server)
	tlsOn := conf.TlsCert != "" && conf.TlsKey != ""
	if tlsOn {
		cert, err := tls.LoadX509KeyPair(conf.TlsCert, conf.TlsKey)
		handle(err)
		grpcServerTls = grpc.NewServer(
			grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
			grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		)
		pb.RegisterKrkStopsServer(grpcServerTls, server)
		grpc_prometheus.Register(grpcServerTls)
	}
	grpc_prometheus.Register(grpcServer)
	if tlsOn {
		grpc_prometheus.Register(grpcServerTls)
	}
	handler := promhttp.Handler()
	http.Handle("/metrics", handler)
	httpServer := http.Server{
		Addr: ":8040",
	}
	lisHttp, err := net.Listen("tcp", ":8040")
	handle(err)
	if tlsOn {
		lisTls, err := net.Listen("tcp", ":9090")
		handle(err)
		go serveGracefully(grpcServerTls, lisTls)
		defer grpcServer.GracefulStop()
	}
	go serveGracefully(grpcServer, lis)
	defer grpcServer.GracefulStop()
	go serveGracefully(&httpServer, lisHttp)
	defer httpServer.Shutdown(context.Background())
	<-stop
}
