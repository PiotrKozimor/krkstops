set -ex
home=/home/Desktop
podman pull docker.io/narciarz96/krkstops-envoy:latest
podman pull docker.io/narciarz96/krkstops
running_pod=$(podman pod ls | grep Running | awk '{print $1}')
new_pod=$(podman pod create -p 9090:9090 -p 6379:6379 -p 8080:8080)
podman create --pod $new_pod docker.io/narciarz96/krkstops-envoy:latest
podman create --pod $new_pod --env-file $home/AIRLY docker.io/narciarz96/krkstops
podman create --pod $new_pod docker.io/redislabs/redisearch
podman pod stop $running_pod
podman pod start $new_pod
podman pod rm -f $running_pod
$home/stopsctl update