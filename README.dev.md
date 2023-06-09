## go generate

1. Fetch stops: `go generate ./pkg/search/score/stops/`
1. Score them: `go generate ./pkg/search/score/`

## go test

Tests require AIRLY_KEY to passed:

```
AIRLY_KEY=foo go test ./...
```