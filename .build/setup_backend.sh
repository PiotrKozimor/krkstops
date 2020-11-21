podman pull docker.io/narciarz96/krkstops-prometheus
podman pull docker.io/narciarz96/krkstops-ttssmonitor
podman pull docker.io/redislabs/redisearch
redis=$(podman create --ip 10.88.0.7 docker.io/redislabs/redisearch)
prometheus=$(podman create --privileged --ip 10.88.0.8 -v /home/piotr/prom:/prometheus docker.io/narciarz96/krkstops-prometheus)
ttssmonitor=$(podman create --ip 10.88.0.10 docker.io/narciarz96/krkstops-ttssmonitor)
podman start $redis $prometheus $ $ttssmonitor
