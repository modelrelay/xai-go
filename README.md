# xai-go

A **gRPC-native Go client for xAI's Grok API**, generated directly from xAI's
official protobuf definitions. It speaks Grok's native gRPC surface with
first-class, strongly-typed Go types and streaming — chat, responses,
embeddings, images, documents, tokenize, models, and auth.

> **Unofficial.** A community client, **not affiliated with, authorized, or
> endorsed by xAI**. "xAI" and "Grok" are trademarks of their respective owner.
> Generated from xAI's public protobuf definitions
> ([`xai-org/xai-proto`](https://github.com/xai-org/xai-proto)) and maintained on
> a best-effort basis by [ModelRelay](https://github.com/modelrelay).
> Licensed under [Apache-2.0](./LICENSE).

### Generated Protos

- `third_party/xai-proto` is tracked as a git submodule.
- Run `make proto` (optionally overriding `BUF`) to regenerate `gen/xai/api/v1`.
- Generation is scoped to `xai/api` (the public inference API). xAI's
  `management_api` and `shared` (billing, analytics) surfaces are intentionally
  excluded — see `PROTO_PATH` in the `Makefile`.
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
	Model: "grok-4.3",
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
	Model: "grok-4.3",
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
	for _, out := range chunk.GetOutputs() {
		fmt.Print(out.GetDelta().GetContent())
	}
}
if outs := acc.Response().GetOutputs(); len(outs) > 0 {
	fmt.Println("\n\nFinal answer:", outs[0].GetMessage().GetContent())
}
```

Prefer a callback? Drain a fresh stream with the high-level helper instead of the manual `Recv` loop:

```go
err = stream.ForEachChunk(ctx, func(chunk *xaiapiv1.GetChatCompletionChunk) error {
	fmt.Printf("\nChunk %s", chunk.GetId())
	return nil
})
if err != nil && err != io.EOF {
	log.Fatal(err)
}
```

### Why gRPC for streaming

Most Grok clients consume the OpenAI-compatible HTTP endpoint, where a stream is
a sequence of Server-Sent Events — `data: {json}\n\n` frames you split,
JSON-decode, and terminate on a `[DONE]` sentinel. This client speaks Grok's
native gRPC surface instead, which changes what streaming feels like from Go:

- **Typed chunks, no frame parsing.** Every `stream.Recv()` returns a
  fully-typed `*GetChatCompletionChunk` — length-prefixed protobuf framing and
  decoding are handled for you. No SSE delimiter splitting, partial-JSON
  reassembly, or `[DONE]` sentinel to special-case.
- **Structured end-of-stream.** A stream ends with a gRPC status: success is a
  clean `io.EOF`, and a mid-stream failure arrives as a typed `status.Code` (via
  HTTP/2 trailers) — not a truncated body or an in-band error event you have to
  sniff for.
- **Real cancellation.** Cancel the `context` and the underlying HTTP/2 stream is
  reset, signaling the server to stop generating; deadlines propagate the same way.
- **One connection, many streams.** HTTP/2 multiplexes concurrent requests over a
  single connection, without per-call connection setup.

For plain one-directional token streaming, SSE works fine and is simpler in the
browser — a chat completion doesn't exercise gRPC's bidirectional streaming. The
win here is consuming the stream as typed, framed, status-terminated messages
from Go, with fewer parsing edge cases.

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
- `client.Batch.*` – create and manage batch jobs and their results.
- `client.Files.*` – upload, list, retrieve, and delete files (with streaming upload/content).
- `client.Video.*` – generate and extend videos, and poll deferred results.

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
	Model:           "grok-4.3",
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

(Default user and user-agent are set in code via `xai.WithDefaultUser` / `xai.WithUserAgent`.)

### Make Targets

- `make proto` – regenerates Go stubs from the pinned proto definitions. Default `BUF` command can be overridden if the buf binary is unavailable.
- `make tidy` – runs `go fmt` across the repo and `go mod tidy`.
- `make ci` – convenience alias for `gofmt` check + `go test ./...` + `buf lint` (also run in CI).
