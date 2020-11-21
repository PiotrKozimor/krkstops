TAG=$1
if [ -z $1 ] 
then
    TAG=latest
fi
echo $TAG
bash .ci/build_ctl.sh
.build/envoy.sh $TAG
.build/krkstops.sh $TAG