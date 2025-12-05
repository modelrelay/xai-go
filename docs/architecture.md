# Grok Go SDK Architecture Plan

## Goals
- Provide an idiomatic Go SDK for Grok/xAI that mirrors the ergonomics of OpenAI’s Responses API where it helps, but is unapologetically gRPC-first (all transports, streaming, and auth are native gRPC concepts).
- Offer both high-level “Responses” workflows and low-level access to every proto-defined service (Chat, Sample, Embedder, Image, Documents, Tokenize, Models, Auth).
- Deliver first-class gRPC streaming support that can emit typed events and accumulate them back into a familiar response structure.
- Keep transport/authentication ergonomics approachable (option pattern, environment defaults) so existing OpenAI Go users can migrate easily even though HTTP compatibility is not a goal.

## Layered Architecture

### Generated Surface
- Keep all protobuf outputs under `gen/xai/apiv1` generated via Buf (`buf.yaml`, `buf.gen.yaml` from `../xai-proto`).
- Generated code (messages + gRPC stubs) is never edited manually; CI runs `buf generate` to keep it in sync with upstream.

### Transport & Core Client
- `internal/transport` owns `grpc.ClientConn` lifecycle, TLS, retries, deadlines, interceptors, and metadata injection (e.g. `Authorization: Bearer ...`, `User-Agent`).
- `internal/raw` wraps each proto client to enforce SDK-wide defaults (request IDs, timeouts, error normalization) before hitting the RPCs.
- `grok.Client` mirrors `openai-go`’s client: constructor wires env defaults, stores shared `option.RequestOption`s, and exposes fields for each high-level service.

### Domain Packages
- `responses`, `chat`, `sample`, `embeddings`, `images`, `documents`, `tokenize`, `models`, `auth` packages translate Go-friendly structs into proto messages and vice versa.
- Shared helper packages (`tools`, `messages`, `usage`, `search`) centralize complex types (tool definitions, content blocks, usage accounting) so callers do not import proto types directly.

-### Responses Layer
+ `responses.Service` exposes `Create`, `CreateStream`, `Retrieve`, `RetrieveStream`, `Delete`, and `StartDeferred` methods, mapping directly onto the Chat RPCs and embracing gRPC streaming primitives (no HTTP compatibility shim required).
+ `CreateStream` wraps `Chat.GetCompletionChunk` streaming RPC to emit typed events: text deltas, reasoning deltas, encrypted content deltas, tool call deltas, usage ticks, finish events, and citations.
+ A streaming accumulator rebuilds the full `GetChatCompletionResponse` as events arrive, giving callers a choice between incremental handling or await-the-end semantics.
+ Tool call events include incremental JSON arguments so applications can dispatch to local handlers and send tool outputs back as new `Message`s. Stored/deferred completion helpers manage IDs end-to-end.

### Supporting Services
- `chat.Service` offers raw access to the entire Chat surface (blocking, streaming, deferred, stored response lifecycle) for advanced flows.
- `sample.Service` wraps `SampleText` and `SampleTextStreaming` for low-level completion sampling (token-level streaming without structured messages).
- `embeddings.Service`, `images.Service`, `documents.Service`, `tokenize.Service`, `models.Service`, and `auth.Service` are thin wrappers over their proto counterparts, sharing transport/auth handling.
- Utilities (e.g. `documents.Connector`) help wire document search or server-side tools (MCP, Web/X search) into Requests’ `tools` definitions.

## Streaming Design
- Use gRPC bidirectional streaming primitives; wrap response stream in `Stream[T]` abstraction similar to `ssestream.Stream`, but backed by gRPC.
- Event types map 1:1 to proto surfaces (`CompletionOutputChunk`, `SamplingUsage` updates, `DebugOutput`), ensuring zero-loss streaming while remaining ergonomic.
- Provide helper functions to convert event streams into channels, iterators, or callback handlers to fit idiomatic Go patterns.

## Deferred & Stored Responses
- Responses service surfaces helper methods for `StartDeferred` / `GetDeferred` / stored response retrieval & deletion, automatically propagating request IDs from Chat responses (`GetChatCompletionResponse.id` / `store_messages`).
- Optional polling helper repeatedly calls `GetDeferredCompletion` until terminal states (`DONE`, `EXPIRED`, `PENDING` timeout).

## Tooling & Search Integration
- Builder APIs construct `Tool` definitions (functions, WebSearch, XSearch, MCP, DocumentSearch) with validation to avoid malformed proto payloads.
- `SearchParameters` helpers make it easy to enable live search, configure result limits, and attach domain restrictions, aligning with Grok’s search features.

## Roadmap Stages
1. **Stage 1 – Bootstrap & Transport**
   - Integrate `xai-proto` via Buf, set up generation scripts, wire module structure, and stand up `grok.Client` plus bare Chat/Sample bindings. ✅
2. **Stage 2 – Responses & Streaming Fundamentals**
   - Build the gRPC-first Responses layer (unary + streaming APIs), event model, accumulator, stored/deferred helpers, and tests around chunk ordering/error propagation.
3. **Stage 3 – Core Service Parity**
   - Ship Embeddings, Images, Documents, Tokenize, Models, and Auth clients so every proto-defined surface is usable from Go early on.
4. **Stage 4 – Tooling & Search Utilities**
   - Provide builders and validators for `Tool` definitions (functions, MCP, searches), search-parameter helpers, and document-search conveniences that integrate tightly with Responses.
5. **Stage 5 – Advanced Orchestration**
   - Add higher-level flows: tool-call dispatch helpers, deferred polling utilities, stored-response lifecycle management, encrypted content handling, and cross-service coordination.
6. **Stage 6 – Hardening & Samples**
   - Integration tests (record/replay or staging), benchmarks, telemetry hooks, and example programs/opinionated guides tailored to gRPC usage.
