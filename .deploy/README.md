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

1. Create an instance with Ubuntu platform image.
2. Install [coreos-installer](https://coreos.github.io/coreos-installer/getting-started/#install-with-cargo).
3. Create and terminate another instance (preserve boot volume, reffered to as second).
4. Attatch second boot volume to first instance using paravirtualized mode.
5. Gnerate butane.ign locally and copy it to first instance: `butane --pretty --strict -d . < butane.yaml > butane.ign`.
6. Install Fedora CodeOS: `./.cargo/bin/coreos-installer install -p metal --ignition-file butane.ign /dev/sdb`.
7. Deatach boot volume from instance.
8.  Create a new instance with second boot volume.

### Hacks

- Allow traffic in [security list](https://cloud.oracle.com/networking/vcns/ocid1.vcn.oc1.eu-frankfurt-1.amaaaaaaxwiopfaaeie4biwvqjs6vps5hejo2lf2hihhvc6nhsqybdbifwgq/security-lists/ocid1.securitylist.oc1.eu-frankfurt-1.aaaaaaaall7lupqxnxxdthyne5faquxtuh2om6jncj2er25rpn3y2toogklq?region=eu-frankfurt-1)
- Follow this: https://stackoverflow.com/a/54810101 :D

## Certificates

```
/bin/podman run --privileged -it --name=certbot --rm --cap-drop all -p 80:80 --volume /var/srv/krkstops/letsencrypt-webroot:/var/lib/letsencrypt:rw,z --volume /var/srv/krkstops/letsencrypt-certs:/etc/letsencrypt:rw,z docker.io/certbot/certbot:arm64v8-latest --standalone --agree-tos -d krkstops.hopto.org -m p1996k@gmail.com certonly
podman secret rm tls-key
cat /var/srv/krkstops/letsencrypt-certs/live/krkstops.hopto.org/privkey.pem | podman secret create tls-key -
podman secret rm tls-cert
cat /var/srv/krkstops/letsencrypt-certs/live/krkstops.hopto.org/fullchain.pem | podman secret create tls-cert -
```