# Deploy krkstops

## Prerequisites

1. `dnf install coreos-installer`
2. Generate `butane.ign`: 
```
butane --pretty --strict -d . < butane.yaml > butane.ign
```

## Azure

Deployment with Fedora Core OS

1. Install azure cli and login.
2. Follow `azure_setup.sh`.

## OCI - Oracle Cloud Infrastructure

### Setup

1. Create an ARM64 instance with Ubuntu platform image.
2. Create and terminate another ARM64 instance within the same availability domain (preserve boot volume, referred to as second).
3. Attach second boot volume to first instance using paravirtualized mode.
4. Generate self-signed certificates in `.deploy` directory: `openssl req -x509 -newkey rsa:4096 -keyout privkey.pem -out fullchain.pem -sha256 -days 365`.
5. Generate butane.ign locally and copy it to first instance: `butane --pretty --strict -d . < butane.yaml > butane.ign`.
6. Install Fedora CodeOS: 
    ```
    docker run --privileged --rm -v /dev:/dev -v /run/udev:/run/udev -v $(pwd)/butane.ign:/data/butane.ign -w /data quay.io/coreos/coreos-installer:release install /dev/sdb -i butane.ign
    ```
7.  Detach boot volume from instance.
8.  Create a new ARM64 instance with second boot volume.

### Hacks

- Allow traffic in [security list](https://cloud.oracle.com/networking/vcns/ocid1.vcn.oc1.eu-frankfurt-1.amaaaaaaxwiopfaaeie4biwvqjs6vps5hejo2lf2hihhvc6nhsqybdbifwgq/security-lists/ocid1.securitylist.oc1.eu-frankfurt-1.aaaaaaaall7lupqxnxxdthyne5faquxtuh2om6jncj2er25rpn3y2toogklq?region=eu-frankfurt-1)
- Follow this: https://stackoverflow.com/a/54810101 :D

## Redisearch

1. Clone [repo](https://github.com/RediSearch/RediSearch)
2. `make setup`
3. `make build`
4. `bin/linux-arm64v8-release/search/redisearch.so` is referenced in `.build/redis.sh` script.

## Certificates

```
/bin/podman run --privileged -it --name=certbot --rm --cap-drop all -p 80:80 --volume /var/srv/krkstops/letsencrypt-webroot:/var/lib/letsencrypt:rw,z --volume /var/srv/krkstops/letsencrypt-certs:/etc/letsencrypt:rw,z docker.io/certbot/certbot:arm64v8-latest --standalone --agree-tos -d krkstops.hopto.org -m p1996k@gmail.com certonly
cp /var/srv/krkstops/letsencrypt-certs/live/krkstops.hopto.org/privkey.pem /etc/krkstops/certs/
cp /var/srv/krkstops/letsencrypt-certs/live/krkstops.hopto.org/fullchain.pem /etc/krkstops/certs/
systemctl restart krkstops.service
```