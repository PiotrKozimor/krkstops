podman pull docker.io/narciarz96/krkstops-envoy:latest
podman pull docker.io/narciarz96/krkstops
podman create --privileged --pod new:krkstops -p 9090:9090 -p 6379:6379 -p 8080:8080 docker.io/narciarz96/krkstops-envoy:latest
podman create --pod krkstops -e AIRLY_KEY=$AIRLY_KEY docker.io/narciarz96/krkstops
podman create --pod krkstops docker.io/redislabs/redisearch
podman pod start krkstops
