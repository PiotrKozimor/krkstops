cd cmd/ttssmonitord
go build 
cd ../..
set -ex
cont=$(buildah from gcr.io/distroless/base-debian10)
buildah copy $cont cmd/ttssmonitord /bin
buildah config --entrypoint '["/bin/ttssmonitord"]' --port 8080 $cont
buildah config --env LOGLEVEL=error $cont
buildah commit --format docker $cont krkstops-ttssmonitor
if [ -z $1 ] 
then
    :
else
    echo "👌 Tag nad push $1"
    buildah tag krkstops-ttssmonitor docker.io/narciarz96/krkstops-ttssmonitor:$1
    buildah push docker.io/narciarz96/krkstops-ttssmonitor:$1
fi
echo "🌝 Tag and push latest"
buildah tag krkstops-ttssmonitor docker.io/narciarz96/krkstops-ttssmonitor:latest
buildah push docker.io/narciarz96/krkstops-ttssmonitor:latest
buildah rm $cont
