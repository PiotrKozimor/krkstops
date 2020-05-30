set -ex
cont=$(buildah from gcr.io/distroless/base-debian10)
buildah copy $cont main/krkstops /bin
buildah config --entrypoint '["/bin/krkstops"]' --port 8080 $cont
buildah commit --format docker $cont krkstops
buildah tag krkstops docker.io/narciarz96/krkstops:latest
buildah push docker.io/narciarz96/krkstops:latest