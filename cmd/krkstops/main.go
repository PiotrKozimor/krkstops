package main

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	_ "time/tzdata"

	"github.com/PiotrKozimor/krkstops"
	"github.com/PiotrKozimor/krkstops/pb"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	server, err := krkstops.NewServer()
	handle(err)

	go func() {
		log.Printf("http rpc server listening on :8082")
		handle(http.ListenAndServe(":8082", server))
	}()

	go func() {
		log.Printf("http metrics server listening on :8083")
		handle(http.ListenAndServe(":8083", krkstops.MetricHandler()))
	}()

	go func() {
		lis, err := net.Listen("tcp", ":8081")
		handle(err)
		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(krkstops.InjectFailure))
		pb.RegisterKrkStopsServer(grpcServer, server)
		log.Printf("grpc staging server listening on :8081")
		handle(grpcServer.Serve(lis))
	}()

	lis, err := net.Listen("tcp", ":8080")
	handle(err)
	grpcServer := grpc.NewServer()
	pb.RegisterKrkStopsServer(grpcServer, server)
	log.Printf("grpc server listening on :8080")
	handle(grpcServer.Serve(lis))
}
