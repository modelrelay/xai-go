# Integration Testing

## Streaming

- Integration tests live in `integration/` and are guarded by the `integration` build tag.
- Run them locally via (disables caching so calls always hit the API):

  ```bash
  XAI_API_KEY=sk-... go test -count=1 -tags=integration ./integration/...
  ```

- Tests expect outbound network access and will skip automatically if `XAI_API_KEY` is unset.

## Replay Harness

Use the replay harness to record real gRPC traffic once and replay it locally without hitting the API.

```go
recorder, err := replay.Open("replay.ndjson", replay.ModeRecord)
if err != nil {
    // handle
}
defer recorder.Close()

client, err := xai.NewClient(ctx,
    xai.WithAPIKey(os.Getenv("XAI_API_KEY")),
    xai.WithDialOptions(recorder.DialOptions()...),
)
```

To replay, load the same file in `ModeReplay`:

```go
replayer, err := replay.Open("replay.ndjson", replay.ModeReplay)
if err != nil {
    // handle
}
defer replayer.Close()

client, err := xai.NewClient(ctx,
    xai.WithAPIKey("replay"),
    xai.WithDialOptions(replayer.DialOptions()...),
)
```
