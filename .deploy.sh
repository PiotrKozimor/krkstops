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
    ssh coreos podman pull docker.io/narciarz96/krkstops:$TAG
    ssh coreos systemctl stop krkstops.service
    ssh coreos podman rm -i krkstops
    ssh coreos podman create \
            --name krkstops \
            --network podman \
            --ip 10.88.0.10 \
            docker.io/narciarz96/krkstops:$TAG
    ssh coreos systemctl start krkstops.service
}

while getopts hk opts; do
   case ${opts} in
    h) help; exit 0;;
    k) krkstops;;
   esac
done




set -xe

