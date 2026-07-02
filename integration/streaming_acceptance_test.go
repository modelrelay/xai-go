//go:build integration

package integration

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/modelrelay/xai-go"
	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/messages"
	"github.com/modelrelay/xai-go/responses"
)

const acceptModel = "grok-4.3"

func acceptReq() *xaiapiv1.GetCompletionsRequest {
	return &xaiapiv1.GetCompletionsRequest{
		Model:    acceptModel,
		Messages: []*xaiapiv1.Message{messages.UserText("Count to five.")},
	}
}

// TestStreamingAcceptance_FreshClientRepeated exercises the exact drain pattern
// shipped in examples/streaming (fresh client + connection each run, stream.Recv
// loop, accumulator), asserting every run yields non-empty content. This is the
// regression guard for the launch-blocking flake where a stream produced zero
// outputs and the example panicked. Each run is bounded by a timeout so a stuck
// stream fails the gate instead of hanging it.
func TestStreamingAcceptance_FreshClientRepeated(t *testing.T) {
	requireKey(t)
	const runs = 20
	for i := 0; i < runs; i++ {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			client, err := xai.NewClient(ctx) // fresh connection each iteration
			if err != nil {
				t.Fatalf("run %d: new client: %v", i, err)
			}
			defer client.Close()
			stream, err := client.Responses.CreateStream(ctx, acceptReq())
			if err != nil {
				t.Fatalf("run %d: CreateStream: %v", i, err)
			}
			// Drain via Recv exactly like examples/streaming: io.EOF ends the
			// stream; any other error is a real gRPC status and must fail (not
			// be swallowed).
			acc := responses.NewAccumulator()
			for {
				chunk, rerr := stream.Recv()
				if errors.Is(rerr, io.EOF) {
					break
				}
				if rerr != nil {
					t.Fatalf("run %d: stream error: %v", i, rerr)
				}
				acc.AddChunk(chunk)
			}
			outs := acc.Response().GetOutputs()
			if len(outs) == 0 || strings.TrimSpace(outs[0].GetMessage().GetContent()) == "" {
				t.Fatalf("run %d: stream produced no content (the flake)", i)
			}
		}()
	}
}
