package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/PiotrKozimor/krkstops"
	"github.com/PiotrKozimor/krkstops/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func serve(s *grpc.Server, l net.Listener) {
	err := s.Serve(l)
	handle(err)
}

type Config struct {
	TlsCert string
	TlsKey  string
}

func main() {
	tlsCert := os.Getenv("TLS_CERT")
	tlsKey := os.Getenv("TLS_KEY")

	server, err := krkstops.NewServer()
	handle(err)

	tlsOn := tlsCert != "" && tlsKey != ""
	if tlsOn {
		cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
		handle(err)
		grpcServerTls := grpc.NewServer(
			grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		)
		pb.RegisterKrkStopsServer(grpcServerTls, server)
		lisTls, err := net.Listen("tcp", ":9090")
		handle(err)
		log.Printf("tls grpc server listening on :9090")
		go serve(grpcServerTls, lisTls)
	}

	go func() {
		Routes()
		var err error
		if tlsOn {
			log.Printf("http server listening on :443")
			err = http.ListenAndServeTLS(":443", tlsCert, tlsKey, nil)
		} else {
			log.Printf("http server listening on :8090")
			err = http.ListenAndServe(":8090", nil)
		}
		handle(err)
	}()

	lis, err := net.Listen("tcp", ":8080")
	handle(err)
	grpcServer := grpc.NewServer()
	pb.RegisterKrkStopsServer(grpcServer, server)
	log.Printf("grpc server listening on :8080")
	serve(grpcServer, lis)
}
