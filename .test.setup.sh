
docker-compose up -d
source .deploy/.env.dev
source .deploy/.env.priv
export AIRLY_KEY
export REDISURI
export OVERRIDE_AIRLY
export OVERRIDE_TTSS_BUS
export OVERRIDE_TTSS_TRAM