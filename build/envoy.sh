set -ex
cont=$(buildah from envoyproxy/envoy:v1.14.1)
buildah copy $cont build/envoy.yaml /etc/envoy/envoy.yaml
buildah copy $cont build/fullchain.pem /etc/pki/fullchain.pem
buildah copy $cont build/privkey.pem /etc/pki/privkey.pem
buildah commit --format docker $cont krkstops-envoy
buildah tag krkstops-envoy docker.io/narciarz96/krkstops-envoy:latest
buildah push docker.io/narciarz96/krkstops-envoy:latest
