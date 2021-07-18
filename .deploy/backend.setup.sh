source .env
podman pull docker.io/redislabs/redisearch:$REDISEARCH_TAG
redis=$(podman create --ip 10.88.0.7 docker.io/redislabs/redisearch:v1.6.15)
podman start $redis
