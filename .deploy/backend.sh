set -ex
home=/home/piotr/krkstops
podman pull docker.io/narciarz96/krkstops:$1
running_app_cont=$(podman container ls | grep krkstops: | awk '{print $1}')
new_app_cont=$(podman create --privileged -p 8080:8080 -p 9090:9090 --ip 10.88.0.5 --env-file $home/.env.prod --env-file $home/.env.priv -v /root/krk-stops-certs:/etc/krkstops:ro docker.io/narciarz96/krkstops:$1)
podman container stop $running_app_cont
podman container start $new_app_cont
podman container rm $running_app_cont
