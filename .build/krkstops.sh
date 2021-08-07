set -ex
cd cmd/krkstops
go build 
cd ../..
cont=$(buildah from gcr.io/distroless/base-debian10)
buildah copy $cont cmd/krkstops/krkstops /bin
buildah config --entrypoint '["/bin/krkstops"]' --port 8080 --port 9090 $cont
buildah commit --format docker $cont krkstops
echo "👌 Tag nad push $1"
buildah tag krkstops docker.io/narciarz96/krkstops:$1
buildah push docker.io/narciarz96/krkstops:$1
