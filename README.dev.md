# Test

Run krkstops:
```
. .test.setup.sh
go run cmd/krkstops/main.go
```
or debug using `krkstops` configuration in VS Code.

In other terminal run tests:
```
. .test.setup.sh
go test ./...
```


# Cache

`bus` - SET holding ids of bus stops

`tram` - SET holding ids of tram stops

`names` - HASH holding mapping id -> name for each stop

`scores` - HASH holding mapping id -> score for each stop

`to.score` - SET holding stops to score

`sug` - key for [suggestions](https://oss.redislabs.com/redisearch/Commands/#suggestions)

## Operations

### Update stops

Please note that data exposed to client will only be affected in last step, when we execute key renaming.

1. Featch bus and tram.
2. Remove `to.score` key. 
3. Add id to `tmp.bus` and/or `tmp.tram`.
4. Filter out non-unique stops (by id, using e.g. map).
5. Add stops names to `tmp.names`.
6. Try to get score for stop from `scores`.
7. Create suggestion in `tmp.sug` for each score. If score was not available, apply temporary score `1` and add stop id to `to.score`.
8. Remove `tmp.` prefix from all keys above (with `RENAME` operation).

### Score stops

For each stop from `to.score`:

1. Pop stop id.
2. Score stop.
3. Save score in `scores`.
4. Create suggestion(s) (old one will be deleted).

### Search stop

1. Execute [`FT.SUGADD`](https://oss.redislabs.com/redisearch/Commands/#ftsugadd) against `sug` key.
2. Get stop name from `names` hash.

# Code gen

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@vlatest
```

# gtfs package

```
cd cmd/gtfs
go run . ../../gtfs/data ../../gtfs/krkstops.gob
go run . stops ../../gtfs/data
```


https://asciiflow.com/#/local/GTFS

```
                 ┌──────────────────┐
          routes │id    PatternText │
                 └▲─────────────────┘
                  │
            ┌─────┴────────────────────┐
      trips │id   route_id   Direction │
            └▲─────────────────────────┘
             │
           ┌─┴─────────────────┐
stop_times │ trip_id   stop_id │
           └───────────────────┘
```