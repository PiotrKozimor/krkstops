set -xe
cont=$(buildah from docker.io/prom/prometheus:v2.25.2)
buildah copy $cont .build/prometheus.yml /etc/prometheus/
buildah commit --format docker $cont krkstops-prometheus
buildah tag krkstops-prometheus docker.io/narciarz96/krkstops-prometheus:latest
buildah push docker.io/narciarz96/krkstops-prometheus:latest
buildah rm $cont