# Streaming Responses Tutorial

This tutorial walks through the recommended streaming pattern using the iterator helper.

## 1. Create a client

```go
ctx := context.Background()
client, err := xai.NewClient(ctx)
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

## 2. Start a streaming request

```go
req := &xaiapiv1.GetCompletionsRequest{
    Model:           "grok-4.3",
    ReasoningEffort: xaiapiv1.ReasoningEffort_EFFORT_LOW.Enum(),
    Messages:        []*xaiapiv1.Message{messages.UserText("Stream a limerick about gRPC.")},
}
stream, err := client.Responses.CreateStream(ctx, req)
if err != nil {
    log.Fatal(err)
}
```

For reasoning models, `ReasoningEffort` is the primary latency lever. Set it to
`EFFORT_NONE`, `EFFORT_LOW`, `EFFORT_MEDIUM`, or `EFFORT_HIGH`; if omitted, the
server defaults to `EFFORT_MEDIUM`.

## 3. Iterate over chunks with cancellation support

```go
acc := responses.NewAccumulator()
it := stream.Iterator(ctx)
for {
    chunk, ok, err := it.Next()
    // Surface a real error before ending on !ok: Next returns ok=false for
    // both a clean EOF and a mid-stream gRPC error.
    if err != nil && !errors.Is(err, io.EOF) {
        log.Fatal(err)
    }
    if !ok {
        break
    }
    acc.AddChunk(chunk)
    for _, out := range chunk.GetOutputs() {
        fmt.Print(out.GetDelta().GetContent())
    }
}
outs := acc.Response().GetOutputs()
if len(outs) == 0 {
    log.Fatal("stream produced no output")
}
fmt.Println("\nFinal:", outs[0].GetMessage().GetContent())
```

## 4. Optional: Use `ForEachChunk`

A stream can only be consumed once, so create a fresh one for the callback helper:

```go
stream, err := client.Responses.CreateStream(ctx, req)
if err != nil {
    log.Fatal(err)
}
if err := stream.ForEachChunk(ctx, func(chunk *xaiapiv1.GetChatCompletionChunk) error {
    fmt.Printf("chunk %s\n", chunk.GetId())
    return nil
}); err != nil && !errors.Is(err, io.EOF) {
    log.Fatal(err)
}
```

See `examples/streaming` for a runnable version.
