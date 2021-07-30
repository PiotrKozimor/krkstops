package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"

	"github.com/PiotrKozimor/krkstops"
	"github.com/PiotrKozimor/krkstops/pb"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", ":8080")
	handle(err)
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	viper.SetDefault("REDIS_URI", "localhost:6379")
	viper.AutomaticEnv()
	redisURI := viper.GetString("REDIS_URI")
	log.Printf("CONFIG: %+v\n", viper.AllSettings())
	server, err := krkstops.NewServer(redisURI)
	handle(err)
	pb.RegisterKrkStopsServer(grpcServer, server)
	grpc_prometheus.Register(grpcServer)
	handler := promhttp.Handler()
	http.Handle("/metrics", handler)
	go func() {
		err := http.ListenAndServe(":8040", nil)
		handle(err)
	}()

	log.Fatal(grpcServer.Serve(lis))
}
