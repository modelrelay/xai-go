## xAI Go SDK

Early-stage Go SDK for xAI's Grok models via gRPC API. Stage 1 focuses on repo scaffolding, proto generation, transport plumbing, and exposing low-level Chat/Sample services.

### Generated Protos

- `third_party/xai-proto` is tracked as a git submodule.
- Run `make proto` (optionally overriding `BUF`) to regenerate `gen/xai/api/v1`.
- Buf managed mode ensures go package paths live under `github.com/modelrelay/xai-go/gen/xai/api/v1`.

### Usage

```go
ctx := context.Background()
client, err := xai.NewClient(ctx, xai.WithAPIKey(os.Getenv("XAI_API_KEY")))
if err != nil {
	log.Fatal(err)
}
defer client.Close()

resp, err := client.Chat.GetCompletion(ctx, &xaiapiv1.GetCompletionsRequest{
	Model: "grok-2-latest",
	Messages: []*xaiapiv1.Message{
		{Role: xaiapiv1.MessageRole_ROLE_USER, Content: []*xaiapiv1.Content{
			{Content: &xaiapiv1.Content_Text{Text: "Hello Grok!"}},
		}},
	},
})
if err != nil {
	log.Fatal(err)
}
fmt.Println(resp.GetOutputs()[0].GetMessage().GetContent())
```

### Streaming Responses

```go
stream, err := client.Responses.CreateStream(ctx, &xaiapiv1.GetCompletionsRequest{
	Model: "grok-2-latest",
	Messages: []*xaiapiv1.Message{
		{Role: xaiapiv1.MessageRole_ROLE_USER, Content: []*xaiapiv1.Content{
			{Content: &xaiapiv1.Content_Text{Text: "Stream something fancy."}},
		}},
	},
})
if err != nil {
	log.Fatal(err)
}

acc := responses.NewAccumulator()
for {
	chunk, err := stream.Recv()
	if err == io.EOF {
		break
	}
	if err != nil {
		log.Fatal(err)
	}
	acc.AddChunk(chunk)
for _, delta := range chunk.GetOutputs() {
	fmt.Print(delta.GetDelta().GetContent())
}
}
full := acc.Response()
fmt.Println("\n\nFinal answer:", full.GetOutputs()[0].GetMessage().GetContent())

// Or use the high-level iterator helper:
if err := stream.ForEachChunk(ctx, func(chunk *xaiapiv1.GetChatCompletionChunk) error {
	fmt.Printf("\nChunk %s", chunk.GetId())
	return nil
}); err != nil && err != io.EOF {
	log.Fatal(err)
}
```

### Deferred Responses

```go
deferred, _ := client.Responses.StartDeferred(ctx, &xaiapiv1.GetCompletionsRequest{/* ... */})
resp, err := client.Responses.PollDeferredCompletion(ctx, deferred.GetRequestId(), 0)
if err != nil {
	log.Fatal(err)
}
fmt.Println(resp.GetOutputs()[0].GetMessage().GetContent())
```

### Additional Services

- `client.Embeddings.Embed` – generate text/image embeddings.
- `client.Images.GenerateImage` – create images from prompts.
- `client.Documents.Search` – query uploaded document collections.
- `client.Tokenize.TokenizeText` – tokenize using Grok models.
- `client.Models.*` – list or inspect available models.
- `client.Auth.GetAPIKeyInfo` – inspect current API key metadata.
- `client.Billing.*` – manage billing info, payment methods, invoices, balances, and spend limits.

Chat requests now include `MaxTurns` to bound agentic tool-calling loops server-side; set `req.MaxTurns` on `GetCompletionsRequest` when you need a hard stop.

### Tool Call Events

```go
tracker := responses.NewToolCallTracker()
stream, _ := client.Responses.CreateStream(ctx, req)
for {
	chunk, err := stream.Recv()
	if err == io.EOF {
		break
	}
	for _, event := range tracker.ConsumeChunk(chunk) {
		if event.ArgumentsDelta != "" {
			fmt.Printf("tool %s args += %s\n", event.CallID, event.ArgumentsDelta)
		}
		if event.Complete {
			fmt.Printf("tool %s ready: %s\n", event.CallID, event.Call.GetFunction().GetArguments())
		}
	}
}
```

### Helper Builders & Utilities

- `messages.UserText/AssistantText/SystemText` help build chat inputs with minimal boilerplate.
- `documents.CollectionSource` + `documents.SearchRequest` simplify document search configuration.
- `tools.FunctionTool`, `tools.WebSearchTool`, etc. build validated tool definitions for Requests payloads.
- `search.Parameters` + `search.WebSource/XSource/RssSource` make it easy to configure live search without manual proto juggling.
- `documents.ToolResponseMessage` turns document search matches into ROLE_TOOL messages for follow-up calls.
- `toolruntime.Registry` converts completed tool events into ROLE_TOOL messages using registered handlers.
- `encrypted.DecryptResponse/DecryptChunk` help you plug in custom decryptors when `use_encrypted_content` is enabled.
- `Responses.RetrieveAndDelete` and `RequireStoredRequests` simplify stored-response lifecycles.
- `config.Config` offers a declarative data map for client instantiation (e.g., load from YAML/JSON and pass to `Config.NewClient`).

```go
fnTool, _ := tools.FunctionTool("lookup_weather", "Fetch weather", map[string]any{
	"type": "object",
	"properties": map[string]any{
		"city": map[string]any{"type": "string"},
	},
})
webSource, _ := search.WebSource(search.WebAllow("example.com"))
params, _ := search.Parameters(
	search.WithMode(xaiapiv1.SearchMode_ON_SEARCH_MODE),
	search.WithSources(webSource),
)
req := &xaiapiv1.GetCompletionsRequest{
	Model:           "grok-2-latest",
	Messages:        []*xaiapiv1.Message{messages.UserText("What's up in Example City?")},
	Tools:           []*xaiapiv1.Tool{fnTool},
	SearchParameters: params,
}

matches, _ := client.Documents.Search(ctx, documents.SearchRequest("Example City history", documents.CollectionSource("city-archive")))
toolMsg, _ := documents.ToolResponseMessage("call_weather_docs", matches.GetMatches(), 3)
_ = toolMsg // send as ROLE_TOOL message when replying to the model.

registry := toolruntime.NewRegistry()
registry.Register("lookup_weather", func(ctx context.Context, fn *xaiapiv1.FunctionCall) (any, error) {
	return map[string]any{"echo": fn.GetArguments()}, nil
})
msg, _ := registry.Handle(ctx, responses.ToolCallEvent{
	CallID:   "call_weather_docs",
	Complete: true,
	Call: &xaiapiv1.ToolCall{
		Tool: &xaiapiv1.ToolCall_Function{Function: &xaiapiv1.FunctionCall{Name: "lookup_weather", Arguments: "{\"city\":\"Example\"}"}}},
	},
})
_ = msg // append ROLE_TOOL message into your conversation
```

### MCP & Server-Side Tool Tips

- Use `tools.WithBearerToken`/`tools.WithAuthorization` to set MCP auth headers.
- Add custom headers via `tools.WithExtraHeader` when MCP servers require proprietary metadata.
- Prefer the search builders above so you never mix mutually exclusive fields (allowed/excluded domains, handles, etc.), keeping requests valid.
- When storing responses (`store_messages=true`) use `responses.RequireStoredRequests` to set `previous_response_id` consistently, and `Responses.RetrieveAndDelete` to clean up stored history once consumed.
- Load multi-environment settings via `config.Config` so deployments can express addresses/API keys in plain data rather than scattering `With*` calls.

### Examples & Integration Tests

- `examples/streaming` demonstrates a streaming chat session using the iterator helper. Run with `go run ./examples/streaming` after setting `XAI_API_KEY`.
- `examples/tool_call` shows how to react to tool-call events, issue document searches, and feed ROLE_TOOL messages back into the conversation.
- Integration tests live under `integration/` and are guarded by the `integration` build tag. Run them with `XAI_API_KEY=... go test -tags=integration ./integration/...`.
- Tutorials:
  - `docs/tutorials/streaming.md` – streaming basics with iterators.
  - `docs/tutorials/tool-doc-search.md` – combines streaming, tool calls, and document search.
  - `docs/tutorials/config.md` – shows how to load settings from JSON/YAML.
- Guides:
  - `docs/guides/responses.md` – full walkthrough covering unary vs streaming, tools, encrypted content, deferred/stored responses.

### Environment Defaults

- `XAI_API_KEY` – required unless `xai.WithAPIKey` is provided.
- `XAI_GRPC_ADDRESS` – optional override for the gRPC endpoint.
- `XAI_DEFAULT_USER` – optional default user identifier for requests.
- `XAI_USER_AGENT` – override the default user agent when needed.

### Make Targets

- `make proto` – regenerates Go stubs from the pinned proto definitions. Default `BUF` command can be overridden if the buf binary is unavailable.
- `make tidy` – runs `go fmt` across the repo and `go mod tidy`.
- `make ci` – convenience alias for `gofmt` check + `go test ./...` + `buf lint` (also run in CI).
