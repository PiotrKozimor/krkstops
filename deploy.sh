#!/bin/bash
set -ex
cd cmd/krkstops
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build .
ssh coreos systemctl stop podman@krkstops.service
scp krkstops coreos:/var/krkstops/krkstops
ssh coreos systemctl start podman@krkstops.service
