cd cmd
for ctl in $(find . -type d)
do  
    if [ "$ctl" == "." ]; then
        echo foobar
    else
        cd $ctl
        echo $ctl
        go build -o ../${ctl}ctl .
        cd ..
    fi
done