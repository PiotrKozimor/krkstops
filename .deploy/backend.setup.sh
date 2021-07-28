redis=$(podman create --ip 10.88.0.7 docker.io/redislabs/redisearch:2.0.10)
podman start $redis
