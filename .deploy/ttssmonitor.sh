running_prom=$(podman container ls | grep krkstops-prometheus: | awk '{print $1}')
running_monitor=$(podman container ls | grep krkstops-ttssmonitor: | awk '{print $1}')
prometheus=$(podman create --privileged --ip 10.88.0.8 -v /home/piotr/prom:/prometheus -v /home/piotr/krkstops/prometheus.yml:/etc/prometheus/prometheus.yml docker.io/prom/prometheus:v2.28.1)
ttssmonitor=$(podman create --ip 10.88.0.10 docker.io/narciarz96/krkstops-ttssmonitor:$1)
podman stop $running_prom $running_monitor
podman start $prometheus $ttssmonitor
