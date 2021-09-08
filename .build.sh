# .build/ctl.sh
TAG=$(git tag --points-at HEAD)
if [ -z $TAG ] 
then
    TAG=$(git rev-parse --short HEAD)
fi
echo "TAG: $TAG"

help () {
	echo "Build"
	echo "	-k krkstops"
	echo "	-t ttssmonitor"
	echo "	-c ctls"
}

while getopts hktc opts; do
   case ${opts} in
    h) help; exit 0;;
    k) .build/krkstops.sh $TAG;;
    t) .build/ttssmonitor.sh $TAG;;
    c) .build/ctl.sh;
   esac
done
