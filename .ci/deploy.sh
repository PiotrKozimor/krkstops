set -e
scp .build/run_backend.sh azure:krkstops
scp .env.prod azure:krkstops/.env
ssh azure sudo bash /home/piotr/krkstops/run_backend.sh $1
rsync -a cmd/*ctl/*ctl azure:krkstops
ssh azure /home/piotr/krkstops/stopsctl update -e 10.88.0.7:6379 -i
