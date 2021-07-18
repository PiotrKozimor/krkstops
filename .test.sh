
help () {
	echo "Test krkstops."
	echo "	-r restart redisearch"
}
while getopts rh opts; do
   case ${opts} in
	  h) help; exit 0;;
	  r) RESTART=true;;
   esac
done

set -e
source .env.dev
source .env.priv
export AIRLY_KEY

redisearch=$(docker container ls | grep redisearch:$REDISEARCH_TAG | awk '{print $1}')
if [ -z $redisearch ] 
then
    if [ $RESTART = "true" ]
    then 
        docker stop $redisearch
        echo "😵 Redisearch is stopped"
    fi
    redisearch=$(docker run -d -p 6380:6379 redislabs/redisearch:$REDISEARCH_TAG)
    echo "🙌 Redisearch is setup"
fi

echo "🚊 Running TTSS test"
cd ttss
go test -v
cd ../airly
echo "🌧️ Running airly test"
go test -v
cd ..
echo "👊 Running krkstops test"
go test -v

echo "🍪 Running airlyctl and ttssctl"
cd cmd/airlyctl
go build
./airlyctl
cd ../ttssctl
go build
./ttssctl deps