TAG=$(git tag --points-at HEAD)
if [ -z $TAG ] 
then
    TAG=$(git rev-parse --short HEAD)
fi
echo "TAG: $TAG"

help () {
	echo "Deploy"
	echo "	-k krkstops"
	echo "	-t ttssmonitor"
	echo "	-c ctls"
}

set -e

deploy_krkstops () {
    scp .deploy/backend.sh azure:krkstops
    scp .env.prod azure:krkstops/.env.prod
    scp .env.priv azure:krkstops/.env.priv
    ssh azure sudo bash /home/piotr/krkstops/backend.sh $TAG
    rsync -a cmd/*ctl/*ctl azure:krkstops
}

deploy_ttssmonitor () {
    scp .deploy/ttssmonitor.sh azure:krkstops
    scp .build/prometheus.yml azure:krkstops
    ssh azure sudo bash /home/piotr/krkstops/ttssmonitor.sh $TAG
}

deploy_ctls () {
    .build/ctl.sh
    rsync -a cmd/*ctl/*ctl azure:krkstops
}

while getopts hktc opts; do
   case ${opts} in
    h) help; exit 0;;
    k) deploy_krkstops;;
    t) deploy_ttssmonitor;;
    c) deploy_ctls;;
   esac
done


