# Integration Testing

## Streaming

- Integration tests live in `integration/` and are guarded by the `integration` build tag.
- Run them locally via (disables caching so calls always hit the API):

  ```bash
  GROK_API_KEY=sk-... go test -count=1 -tags=integration ./integration/...
  ```

- Tests expect outbound network access and will skip automatically if `GROK_API_KEY` is unset.

## Future Replay Harness (TODO)

- Record gRPC responses in staging and replay locally for deterministic CI.
- Track open work in Stage 7.
