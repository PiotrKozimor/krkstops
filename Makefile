deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protos: deps
	protoc pb/krk-stops.proto --go_out=.
	protoc pkg/gtfs/realtimepb/gtfs-realtime.proto --go_out=. --go_opt paths=source_relative --go_opt=Mpkg/gtfs/realtimepb/gtfs-realtime.proto=go.cozymore.dev/krkstops/gtfs/realtimepb