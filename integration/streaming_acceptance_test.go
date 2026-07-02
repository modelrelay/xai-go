//go:build integration

package integration

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

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

// TestStreamingAcceptance_ExamplePattern mirrors examples/streaming EXACTLY
// (fresh client, Iterator, manual AddChunk), across many fresh connections,
// and reports the terminal (ok,err) whenever a run yields 0 outputs — the
// case where the example panics. Captures the error the example's loop swallows.
func TestStreamingAcceptance_ExamplePattern(t *testing.T) {
	requireKey(t)
	const runs = 30
	fails := 0
	for i := 0; i < runs; i++ {
		func() {
			ctx := context.Background()
			client, err := xai.NewClient(ctx)
			if err != nil {
				t.Fatalf("run %d: new client: %v", i, err)
			}
			defer client.Close()
			stream, err := client.Responses.CreateStream(ctx, acceptReq())
			if err != nil {
				t.Fatalf("run %d: CreateStream: %v", i, err)
			}
			acc := responses.NewAccumulator()
			it := stream.Iterator(ctx)
			n := 0
			var termOK bool
			var termErr error
			for {
				chunk, ok, nerr := it.Next()
				termOK, termErr = ok, nerr
				if errors.Is(nerr, io.EOF) || !ok {
					break
				}
				if nerr != nil {
					break
				}
				n++
				acc.AddChunk(chunk)
			}
			outs := acc.Response().GetOutputs()
			if len(outs) == 0 || strings.TrimSpace(outs[0].GetMessage().GetContent()) == "" {
				fails++
				t.Errorf("run %d FLAKE: chunks=%d outputs=%d term(ok=%v err=%v)", i, n, len(outs), termOK, termErr)
			}
		}()
	}
	if fails > 0 {
		t.Fatalf("%d/%d runs hit the 0-output flake", fails, runs)
	}
}
