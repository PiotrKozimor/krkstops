podman pull docker.io/narciarz96/krkstops-prometheus
podman pull docker.io/redislabs/redisearch
redis=$(podman create --ip 10.88.0.7 docker.io/redislabs/redisearch)
prometheus=$(podman create --ip 10.88.0.8 docker.io/narciarz96/krkstops-prometheus)
podman start $redis $prometheus
