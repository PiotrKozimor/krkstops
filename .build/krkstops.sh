set -ex
cd cmd/krkstops
CGO_ENABLED=0 go build .
cd ../..
cont=$(buildah from scratch)
buildah copy $cont cmd/krkstops/krkstops /bin/krkstops
buildah config --entrypoint '["/bin/krkstops"]' --port 8080 --port 9090 $cont
buildah commit $cont krkstops
echo "👌 Tag nad push $1"
buildah tag krkstops docker.io/narciarz96/krkstops:$1
buildah push docker.io/narciarz96/krkstops:$1
