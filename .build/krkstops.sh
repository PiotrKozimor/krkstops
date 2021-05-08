set -ex
cont=$(buildah from gcr.io/distroless/base-debian10)
buildah copy $cont cmd/krkstops/krkstops /bin
buildah config --entrypoint '["/bin/krkstops"]' --port 8080 $cont
buildah commit --format docker $cont krkstops
if [ -z $1 ] 
then
    :
else
    echo "👌 Tag nad push $1"
    buildah tag krkstops docker.io/narciarz96/krkstops:$1
    buildah push docker.io/narciarz96/krkstops:$1
fi
echo "🌝 Tag and push latest"
buildah tag krkstops docker.io/narciarz96/krkstops:latest
buildah push docker.io/narciarz96/krkstops:latest
buildah rm $cont
