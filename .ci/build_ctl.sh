cd cmd
for ctl in $(find . -type d)
do  
    if [ "$ctl" == "." ]; then
        echo skip
    else
        cd $ctl
        go build .
        echo $ctl
        cd ..
    fi
done
