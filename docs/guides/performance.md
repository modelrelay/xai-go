# Performance Playbook

## Benchmarks

- `go test -bench . ./responses` runs the accumulator benchmark.
- Track improvements when optimizing streaming aggregation.

## gRPC Keepalive

- Set `grpc.WithKeepaliveParams` via `xai.WithDialOptions` to keep long-lived streams healthy.

```go
kac := keepalive.ClientParameters{Time: 30 * time.Second, Timeout: 10 * time.Second}
client, err := xai.NewClient(ctx, xai.WithDialOptions(grpc.WithKeepaliveParams(kac)))
if err != nil {
    log.Fatal(err)
}
```

## Reuse Accumulators

- For high-volume streaming, reuse an `Accumulator` across related requests to reduce allocations.
- Call `acc.Reset()` or instantiate from a sync.Pool.

## Streaming Best Practices

- Use the iterator with context cancellation to cleanly abort streams.
- Decrypt chunks inline to avoid extra passes over the data.
- Avoid copying entire chunks unless necessary; process deltas in-place.
