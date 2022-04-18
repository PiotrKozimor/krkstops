set -ex
cd cmd/ttssmonitord
GOARCH=arm64 CGO_ENABLED=0 go build .
cd ../..
cont=$(buildah from --arch arm64 scratch)
buildah copy $cont cmd/ttssmonitord/ttssmonitord /bin/ttssmonitord
buildah config --entrypoint '["/bin/ttssmonitord"]' --port 8080 $cont
buildah config --env LOGLEVEL=error $cont
buildah commit $cont krkstops-ttssmonitor
echo "👌 Tag nad push $1"
buildah tag krkstops-ttssmonitor docker.io/narciarz96/krkstops-ttssmonitor:$1
buildah push docker.io/narciarz96/krkstops-ttssmonitor:$1
