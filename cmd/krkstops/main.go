package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"

	"github.com/PiotrKozimor/krkstops"
	"github.com/PiotrKozimor/krkstops/airly"
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/ttss"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	viper.SetDefault("REDIS_URI", "localhost:6379")
	viper.AutomaticEnv()
	redisURI := viper.GetString("REDIS_URI")
	log.Printf("CONFIG: %+v\n", viper.AllSettings())
	cache, err := krkstops.NewCache(redisURI, krkstops.SUG)
	if err != nil {
		log.Fatal(err)
	}
	server := krkstops.KrkStopsServer{
		C:     cache,
		Airly: airly.Api,
		Ttss:  ttss.KrkStopsEndpoints,
	}
	pb.RegisterKrkStopsServer(grpcServer, &server)
	grpc_prometheus.Register(grpcServer)
	handler := promhttp.Handler()
	http.Handle("/metrics", handler)
	go func() {
		err := http.ListenAndServe(":8040", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Fatal(grpcServer.Serve(lis))
}
