set -ex
cd cmd/ttssmonitord
go build 
cd ../..
cont=$(buildah from gcr.io/distroless/base-debian10)
buildah copy $cont cmd/ttssmonitord /bin
buildah config --entrypoint '["/bin/ttssmonitord"]' --port 8080 $cont
buildah config --env LOGLEVEL=error $cont
buildah commit --format docker $cont krkstops-ttssmonitor
echo "👌 Tag nad push $1"
buildah tag krkstops-ttssmonitor docker.io/narciarz96/krkstops-ttssmonitor:$1
buildah push docker.io/narciarz96/krkstops-ttssmonitor:$1
