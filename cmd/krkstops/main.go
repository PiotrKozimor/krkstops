package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/PiotrKozimor/krkstops"
	"github.com/PiotrKozimor/krkstops/airly"
	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/ttss"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v8"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var ctx = context.Background()

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
	redisearchEndpoint := os.Getenv("REDISEARCH_ENDPOINT")
	log.Printf("REDISEARCH_ENDPOINT: %s\n", redisearchEndpoint)
	server := krkstops.KrkStopsServer{
		C: krkstops.Clients{
			RedisAutocompleter: redisearch.NewAutocompleter(redisearchEndpoint, "search-stops"),
			Redis:              redis.NewClient(&redis.Options{Addr: redisearchEndpoint})},
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
