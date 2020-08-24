cd ctl
for ctl in $(find . -type f -name "*ctl.go")
do 
    go build $ctl
done
cd ../main
go build krkstops.go
build/envoy.sh
build/krkstops.sh