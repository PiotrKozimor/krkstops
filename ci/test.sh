golint -set_exit_status main/krkstops.go
staticcheck main/krkstops.go
go test ./krkstops
cd ctl
for ctl in $(find . -type f -name "*ctl.go")
do 
    go build $ctl
done
cd ../main
go build krkstops.go