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
    scp .deploy/.tags.env azure:/etc/krkstops/tags.env
    scp .env.prod azure:krkstops/.env.prod
    scp .env.priv azure:krkstops/.env.priv
    ssh azure sudo systemctl restart krkstops.service
}

deploy_ttssmonitor () {
    scp .deploy/.tags.env azure:/etc/krkstops/tags.env
    ssh azure sudo systemctl restart ttssmonitor.service
}

deploy_ctls () {
    .build/ctl.sh
    rsync -a cmd/*ctl/*ctl azure:
}

while getopts hktc opts; do
   case ${opts} in
    h) help; exit 0;;
    k) deploy_krkstops;;
    t) deploy_ttssmonitor;;
    c) deploy_ctls;;
   esac
done


