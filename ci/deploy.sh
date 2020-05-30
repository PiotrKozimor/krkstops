set -e
# rsync ctl/* azure:
# build/envoy.sh
# build/krkstops.sh
scp build/run_backend.sh azure:
scp AIRLY azure:
ssh azure sudo bash run_backend.sh