cd cmd
for ctl in $(ls -d -- */) #find . -type d)
do  
    if [ "ctl" == "." ]; then
        echo foo
    else
        cd $ctl
        echo $ctl
        # echo $(pwd)

        go build -o ../${ctl}ctl .
        cd ..
    fi
done
# cd ../main
# go build krkstops.go
# cd ..
# build/envoy.sh
# build/krkstops.sh