## go generate

1. Fetch stops: `go generate ./pkg/search/score/stops/`
1. Score them: `go generate ./pkg/search/score/`

## go test

Pass the actual AIRLY_KEY in tests:

```
AIRLY_KEY=foo go test ./...
```