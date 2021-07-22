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
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v8"
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
	viper.SetDefault("REDISEARCH_URI", "localhost:6379")
	viper.AutomaticEnv()
	redisearchUri := viper.GetString("REDISEARCH_URI")
	log.Printf("CONFIG: %+v\n", viper.AllSettings())
	server := krkstops.KrkStopsServer{
		C: krkstops.Clients{
			RedisAutocompleter: redisearch.NewAutocompleter(redisearchUri, "search-stops"),
			Redis:              redis.NewClient(&redis.Options{Addr: redisearchUri})},
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
