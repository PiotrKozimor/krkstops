podman run --interactive --rm -v $(pwd):/var/krkstops quay.io/coreos/butane:release \
       --strict --files-dir /var/krkstops < butane.yaml > butane.ign