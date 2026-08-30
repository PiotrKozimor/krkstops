package krkstops

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"go.cozymore.dev/krkstops/pb"
	"google.golang.org/protobuf/proto"
)

func (s *KrkStopsServer) serveHTTP(rw http.ResponseWriter, r *http.Request) error {

	lookup := map[string]struct {
		alloc func() proto.Message
		call  func(context.Context, any) (proto.Message, error)
	}{
		"/3": {func() proto.Message { return &pb.SearchStopsRequest{} }, func(ctx context.Context, a any) (proto.Message, error) {
			return s.SearchStops(ctx, a.(*pb.SearchStopsRequest))
		}},
		"/4": {func() proto.Message { return &pb.GetDeparturesRequest{} }, func(ctx context.Context, a any) (proto.Message, error) {
			return s.GetDepartures(ctx, a.(*pb.GetDeparturesRequest))
		}},
	}

	impl, ok := lookup[r.URL.Path]
	if !ok {
		return fmt.Errorf("unknown path: %s", r.URL.Path)
	}
	b, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}
	request := impl.alloc()
	err = proto.Unmarshal(b, request)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	response, err := impl.call(r.Context(), request)
	if err != nil {
		return fmt.Errorf("call impl: %w", err)
	}
	b, err = proto.Marshal(response)
	if err != nil {
		return fmt.Errorf("marshall: %w", err)
	}

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") || slices.Contains(r.Header.Values("Accept-Encoding"), "gzip") {
		if len(b) > 500 {
			rw.Header().Set("Content-Encoding", "gzip")
			buf := bytes.Buffer{}
			w := gzip.NewWriter(&buf)
			_, err1 := w.Write(b)
			err2 := w.Close()
			_, err3 := rw.Write(buf.Bytes())
			return errors.Join(err1, err2, err3)
		}
	}
	_, err = rw.Write(b)
	return err
}

func (s *KrkStopsServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	err := s.serveHTTP(rw, r)
	if err != nil {
		log.Print("request failed: ", err)
		http.Error(rw, err.Error(), 500)
		requests.With(prometheus.Labels{"status_code": "500", "path": r.URL.Path}).Inc()
	} else {
		requests.With(prometheus.Labels{"status_code": "200", "path": r.URL.Path}).Inc()
	}
}
