set -xe
cont=$(buildah from prom/prometheus)
buildah copy $cont ./prometheus.yml /etc/prometheus/
buildah commit --format docker $cont krkstops-prometheus
buildah tag krkstops-prometheus docker.io/narciarz96/krkstops-prometheus:latest
buildah push docker.io/narciarz96/krkstops-prometheus:latest
buildah rm $cont