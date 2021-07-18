.build/ctl.sh
TAG$=$(git tag --points-at HEAD)
if [ -z $TAG ]
    TAG=$(git rev-parse --short HEAD)
fi
.build/krkstops.sh $TAG