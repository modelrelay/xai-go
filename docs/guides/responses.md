# Responses Guide

This guide covers the full lifecycle of Grok responses: unary completions, streaming, tool calls, encrypted content, and stored responses.

## Unary vs Streaming

- **Unary (`Responses.Create`)** — simpler, blocks until completion, ideal for short answers or batch workloads.
- **Streaming (`Responses.CreateStream`)** — emits incremental chunks, lets you display tokens as they arrive, required for tool orchestration.

## Streaming Pattern

```go
stream, _ := client.Responses.CreateStream(ctx, req)
acc := responses.NewAccumulator()
tracker := responses.NewToolCallTracker()
registry := toolruntime.NewRegistry()

registry.Register("lookup_docs", func(ctx context.Context, fn *xaiapiv1.FunctionCall) (any, error) {
    resp, err := client.Documents.Search(ctx, documents.SearchRequest(fn.GetArguments(), documents.CollectionSource("docs")))
    if err != nil {
        return nil, err
    }
    return resp.GetMatches(), nil
})

if err := stream.ForEachChunk(ctx, func(chunk *xaiapiv1.GetChatCompletionChunk) error {
    acc.AddChunk(chunk)
    for _, event := range tracker.ConsumeChunk(chunk) {
        if !event.Complete {
            continue
        }
        msg, err := registry.Handle(ctx, event)
        if err != nil {
            return err
        }
        // append ROLE_TOOL message back into conversation
        _ = msg
    }
    return nil
}); err != nil {
    log.Fatal(err)
}

final := acc.Response()
fmt.Println(final.GetOutputs()[0].GetMessage().GetContent())
```

## Encrypted Content

When `use_encrypted_content=true`, decrypt payloads before consuming them:

```go
decode := func(cipher string) (string, error) {
    return decrypt(cipher), nil
}

if err := encrypted.DecryptChunk(chunk, decode); err != nil {
    log.Fatal(err)
}
```

Apply `DecryptResponse` to unary completions.

## Stored Responses

```go
req := &xaiapiv1.GetCompletionsRequest{StoreMessages: true}
responses.RequireStoredRequests(req, previousResponseID)

resp, _ := client.Responses.Create(ctx, req)
stored, _ := client.Responses.Retrieve(ctx, resp.GetId())
client.Responses.Delete(ctx, resp.GetId())
```

Use `RetrieveAndDelete` to fetch and clean up in one call.

## Deferred Completions

```go
def, _ := client.Responses.StartDeferred(ctx, req)
full, err := client.Responses.PollDeferredCompletion(ctx, def.GetRequestId(), 0)
if err != nil {
    log.Fatal(err)
}
```

## Putting It Together

1. Create request with tools/search parameters.
2. Start stream.
3. Use iterator + tracker + registry and optional decryptor.
4. Append ROLE_TOOL messages to conversation state.
5. Inspect accumulator for final response and optionally store/delete.

See also `examples/tool_call` for a runnable example.
