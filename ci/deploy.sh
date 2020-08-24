set -e
rsync ctl/* azure:krkstops
scp build/run_backend.sh azure:krkstops
scp AIRLY azure:krkstops
ssh azure sudo bash krkstops/run_backend.sh
ssh azure ./krkstops/stopsctl update