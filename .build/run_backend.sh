set -ex
home=/home/piotr/krkstops
podman pull docker.io/narciarz96/krkstops-envoy:$1
podman pull docker.io/narciarz96/krkstops:$1
running_envoy_cont=$(podman container ls | grep krkstops-envoy: | awk '{print $1}')
running_app_cont=$(podman container ls | grep krkstops: | awk '{print $1}')
new_app_cont=$(podman create -p 8080:8080 --ip 10.88.0.5 --env-file $home/KRKSTOPS docker.io/narciarz96/krkstops:$1)
new_envoy_cont=$(podman create --privileged -p 9090:9090 --ip 10.88.0.6 -v /root/krk-stops-certs/:/etc/pki/:ro docker.io/narciarz96/krkstops-envoy:$1)
podman container stop $running_envoy_cont
podman container stop $running_app_cont
podman container start $new_app_cont $new_envoy_cont
podman container rm $running_envoy_cont $running_app_cont
