# Grok Go SDK Roadmap

## Stage 1 – Bootstrap & Transport
- [x] Pin `xai-proto` as a Git submodule or Go workspace reference and add Buf config to this repo.
- [x] Wire `make proto` (or `go:generate`) that runs `buf generate` and writes outputs into `gen/xai/apiv1`.
- [x] Initialize Go module (`grok-go-sdk`) and bring in base deps (`google.golang.org/grpc`, Buf runtime, option utilities).
- [x] Implement `internal/transport` for `grpc.ClientConn` creation, auth metadata, retries, and logging interceptors.
- [x] Implement `internal/raw` wrappers for Chat/Sample/etc. that apply shared request defaults.
- [x] Create `grok.Client`, hook up env-based defaults (API key, base URL), and expose basic `ChatService`.

## Stage 2 – Responses & Streaming Fundamentals
- [x] Define gRPC-native `responses` types and interfaces that sit directly atop Chat RPCs (no HTTP compatibility shims).
- [x] Implement unary `Create`, `Retrieve`, `Delete`, `StartDeferred`, `GetDeferred` flows plus stored response helpers.
- [x] Implement `CreateStream`/`RetrieveStream`, including event decoding for content, reasoning, tool calls, encrypted payloads, and usage deltas.
- [x] Build a streaming accumulator that reconstructs `GetChatCompletionResponse` and exposes channel/callback helpers.
- [x] Add tool-call event helpers (incremental JSON args, completion detection, error surfacing).
- [x] Expose the streaming accumulator as a pure reducer (`func(state, chunk) state`) so functional consumers can compose their own pipelines.
- [x] Write unit tests covering accumulator behavior and chunk merging correctness.

## Stage 3 – Core Service Parity
- [x] Implement Embeddings, Images, Documents, Tokenize, Models, and Auth services with idiomatic request/response translators.
- [x] Ensure each service shares transport/auth defaults via `internal/raw` and exposes both unary + streaming (if available) entry points.
- [x] Add shared helper packages for content blocks, usage structs, and document search sources.
- [x] Provide service-level examples/tests demonstrating basic interactions.

## Stage 4 – Tooling & Search Utilities
- [x] Implement builders/validators for `Tool` definitions (functions, WebSearch, XSearch, MCP, DocumentSearch, CollectionsSearch).
- [x] Ship `SearchParameters` helper APIs (live search toggles, domain restrictions, MCP extras) with validation.
- [x] Provide convenience helpers for composing Documents search + tool responses, plus wiring MCP auth headers.
- [x] Document best practices for tool definition reuse and server-side search toggles.

## Stage 5 – Advanced Orchestration
- [x] Provide tool orchestration helpers (register local handlers, dispatch tool calls, feed results back into Chat/Responses flows).
- [x] Add deferred-completion polling utilities with context cancellation and backoff controls.
- [x] Support encrypted content helpers (decrypt stubs, inspection toggles) and stored response lifecycle utilities (list, retrieve, delete).
- [x] Offer guidance/examples for combining Documents search, live search, and tool calls inside Responses.
- [x] Introduce a declarative config/data map (no global env reliance) for composing transports in a data-first style.

## Stage 6 – Hardening & Samples
- [x] Add integration tests (record/replay harness or staging endpoint) covering unary + streaming flows.
- [x] Build benchmarks for streaming throughput and accumulator overhead.
- [x] Port key usage examples (chat completion, responses, tool calling, documents search) to demonstrate gRPC-first flows.
- [x] Provide README/API docs describing setup, env vars, migration tips, and best practices.
- [x] Configure CI (lint, go test, Buf format/lint, generated-code verification).

## Stage 7 – Documentation & Guides
- [x] Author an end-to-end Responses guide (unary vs streaming, tool calls, encrypted content, stored responses).
- [x] Write quickstart snippets for Embeddings, Images, Documents, Tokenize, Models, and Auth.
- [x] Expand the tool orchestration deep dive with guidance on feeding ROLE_TOOL messages back into conversations.
- [x] Document integration testing best practices (running `-tags integration`, local creds, future replay harness).
- [x] Add a performance playbook section covering benchmarks, keepalive tuning, and accumulator reuse tips.
