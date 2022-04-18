cont=$(buildah from --arch=arm64 docker.io/redis:6.2.6)
buildah copy $cont redisearch.so /usr/local/lib/
buildah config --entrypoint '["/usr/local/bin/redis-server", "--loadmodule", "/usr/local/lib/redisearch.so"]' --cmd "" $cont
buildah commit $cont redisearch
echo "👌 Tag nad push $1"
buildah tag redisearch docker.io/narciarz96/redisearch:$1
buildah push docker.io/narciarz96/redisearch:$1
