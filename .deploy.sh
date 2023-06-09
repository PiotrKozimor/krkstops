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
            -p 8080:8080 \
            -p 9090:9090 \
            -p 443:443 \
            --secret airly-key,type=env,target=AIRLY_KEY \
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

