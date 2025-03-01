#!/bin/bash
set -e

TAG=$(git tag --points-at HEAD)
if [ -z $TAG ] 
then
    TAG=$(git rev-parse --short HEAD)
fi
echo "TAG: $TAG"

help () {
	echo "Deploy"
	echo "	-k krkstops"
}

krkstops () {
    ssh coreos sudo podman pull docker.io/narciarz96/krkstops:$TAG
    ssh coreos sudo systemctl stop krkstops.service
    ssh coreos sudo podman rm -i krkstops
    ssh coreos sudo podman create \
            --name krkstops \
            --network podman \
            --ip 10.88.0.10 \
            -p 9090:9090 \
            --secret tls-cert \
            --secret tls-key \
            --env TLS_CERT=/run/secrets/tls-cert \
            --env TLS_KEY=/run/secrets/tls-key \
            docker.io/narciarz96/krkstops:$TAG
    ssh coreos sudo systemctl start krkstops.service
}

while getopts hk opts; do
   case ${opts} in
    h) help; exit 0;;
    k) krkstops;;
   esac
done




set -xe

