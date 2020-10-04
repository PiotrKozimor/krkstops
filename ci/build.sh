bash ci/build_ctl.sh
cd ../main
go build krkstops.go
cd ..
build/envoy.sh
build/krkstops.sh