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
    if errors.Is(err, io.EOF) || !ok {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    acc.AddChunk(chunk)
    for _, out := range chunk.GetOutputs() {
        fmt.Print(out.GetDelta().GetContent())
    }
}
fmt.Println("\nFinal:", acc.Response().GetOutputs()[0].GetMessage().GetContent())
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
