source .env.dev
source .env.priv
export REDISEARCH_ENDPOINT
export AIRLY_KEY
go run cmd/krkstops/main.go