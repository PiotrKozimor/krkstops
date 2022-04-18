#!/bin/bash
TAG=$(git tag --points-at HEAD)
if [ -z $TAG ] 
then
    TAG=$(git rev-parse --short HEAD)
fi
echo "TAG: $TAG"

echo 

help () {
	echo "Deploy"
	echo "	-a app"
	echo "	-c ctls"
}

set -e

echo TAG=$TAG > .deploy/.tags.env

deploy_app () {
    scp .deploy/.tags.env azure:/etc/krkstops/tags.env
    scp .deploy/.env.prod azure:/etc/krkstops/.env.prod
    scp .deploy/.env.priv azure:/etc/krkstops/.env.priv
    ssh azure sudo systemctl restart krkstops.service
    ssh azure sudo systemctl restart ttssmonitor.service
}

deploy_ctls () {
    rsync --progress -a cmd/*ctl/*ctl coreos:
}

while getopts hac opts; do
   case ${opts} in
    h) help; exit 0;;
    a) deploy_app;;
    c) deploy_ctls;;
   esac
done


