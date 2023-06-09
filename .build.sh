#!/bin/bash
set -e

TAG=$(git tag --points-at HEAD)
if [ -z $TAG ] 
then
    TAG=$(git rev-parse --short HEAD)
fi
echo "TAG: $TAG"

help () {
	echo "Build"
	echo "	-k krkstops"
}

krkstops () {
    cd cmd/krkstops
    GOARCH=arm64 CGO_ENABLED=0 go build .
    cd ../..
    cont=$(buildah from --arch arm64 scratch)
    buildah copy $cont cmd/krkstops/krkstops /bin/krkstops
    buildah add $cont https://curl.se/ca/cacert.pem /etc/pki/tls/certs/ca-bundle.crt
    buildah config --entrypoint '["/bin/krkstops"]' --port 8080 --port 9090 $cont
    buildah commit $cont krkstops
    echo "👌 Tag nad push $TAG"
    buildah tag krkstops docker.io/narciarz96/krkstops:$TAG
    buildah push docker.io/narciarz96/krkstops:$TAG
}

while getopts hk opts; do
   case ${opts} in
    h) help; exit 0;;
    k) krkstops;;
   esac
done


