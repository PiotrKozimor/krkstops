set -e
rsync cmd/* azure:krkstops
scp build/run_backend.sh azure:krkstops
scp KRKSTOPS azure:krkstops
ssh azure sudo bash /home/piotr/krkstops/run_backend.sh
ssh azure /home/piotr/krkstops/stopsctl update -e 10.88.0.7:6379 -i
