source .deploy/.env.dev
source .deploy/.env.priv
export REDIS_URI
export AIRLY_KEY
go run cmd/krkstops/main.go