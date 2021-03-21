set -e
sudo systemctl start redis

source KRKSTOPS
export AIRLY_KEY
cd ttss
go test -v
cd ../airly
go test -v
cd ../stops
go test -v
cd ../cache
go test -v

sudo systemctl stop redis