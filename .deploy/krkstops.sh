podman create \
  --name krkstops \
  -v /var/krkstops:/app:z \
  --workdir /app \
  --entrypoint /app/krkstops \
  --network podman \
  --ip 10.88.0.10 \
  docker.io/library/alpine:latest