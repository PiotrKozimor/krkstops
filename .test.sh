
help () {
	echo "Test krkstops."
	echo "	-r restart redisearch"
}

restart () {
    docker-compose down
    docker-compose up -d
}

while getopts rh opts; do
   case ${opts} in
	  h) help; exit 0;;
	  r) restart;;
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

go test -v ./...
