source .deploy/.env.dev
source .deploy/.env.priv
export REDISURI
export AIRLY_KEY
go run cmd/krkstops/main.go