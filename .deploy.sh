#!/bin/bash
set -e
cd cmd/krkstops
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build .
ssh coreos systemctl stop krkstops.service
scp krkstops coreos:/var/krkstops/krkstops
ssh coreos systemctl start krkstops.service
