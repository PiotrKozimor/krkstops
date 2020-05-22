podman create --privileged --pod new:krkstops1 -v /root/envoy.yaml:/etc/envoy/envoy.yaml -p 9090:9090 -p 6379:6379 -p 8080:8080 envoyproxy/envoy:v1.14.1
podman create --pod krkstops1 -e AIRLY_KEY=$AIRLY_KEY docker.io/narciarz96/krkstops
podman create --pod krkstops1 docker.io/redislabs/redisearch