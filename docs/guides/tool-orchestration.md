# Tool Orchestration Deep Dive

This guide focuses on managing ROLE_TOOL messages and integrating tool outputs back into the conversation.

## Components

- `responses.ToolCallTracker` — parses streaming chunks and emits `ToolCallEvent`s.
- `toolruntime.Registry` — maps function names to handlers.
- `documents.ToolResponseMessage` — helper to format document search output as ROLE_TOOL content.

## Workflow

```go
tracker := responses.NewToolCallTracker()
registry := toolruntime.NewRegistry()

registry.Register("lookup_docs", func(ctx context.Context, fn *xaiapiv1.FunctionCall) (any, error) {
    matches, err := client.Documents.Search(ctx, documents.SearchRequest(fn.GetArguments(), documents.CollectionSource("docs")))
    if err != nil {
        return nil, err
    }
    return matches.GetMatches(), nil
})

stream, err := client.Responses.CreateStream(ctx, reqWithTools)
if err != nil {
    log.Fatal(err)
}
acc := responses.NewAccumulator()

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
        // Append msg to your conversation history and continue streaming
    }
    return nil
}); err != nil {
    log.Fatal(err)
}
```

### Feeding ROLE_TOOL Messages Back

1. Maintain your own conversation slice (same type used in requests).
2. When `registry.Handle` returns a message, append it to the slice.
3. Resume streaming or send a follow-up request with the updated messages.

### Error Handling Tips

- If a handler fails, consider returning a ROLE_TOOL message describing the failure instead of aborting the stream.
- Use context deadlines to keep tool handlers bounded.
- Combine `documents.ToolResponseMessage` with custom metadata when you need structured outputs.
