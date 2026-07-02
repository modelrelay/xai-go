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
    Model:    "grok-4.3",
    Messages: []*xaiapiv1.Message{messages.UserText("Stream a limerick about gRPC.")},
}
stream, err := client.Responses.CreateStream(ctx, req)
if err != nil {
    log.Fatal(err)
}
```

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

```go
if err := stream.ForEachChunk(ctx, func(chunk *xaiapiv1.GetChatCompletionChunk) error {
    fmt.Printf("chunk %s\n", chunk.GetId())
    return nil
}); err != nil {
    log.Fatal(err)
}
```

See `examples/streaming` for a runnable version.
