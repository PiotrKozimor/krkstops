#!/bin/bash
protoc \
    pb/krk-stops.proto \
     --go_out=. \
     --go-grpc_out=.