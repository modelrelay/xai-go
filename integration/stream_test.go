//go:build integration

package integration

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/modelrelay/xai-go"
	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/messages"
	"github.com/modelrelay/xai-go/responses"
)

func TestStreamingIntegration(t *testing.T) {
	if os.Getenv("XAI_API_KEY") == "" {
		t.Skip("set XAI_API_KEY to run integration test")
	}
	ctx := context.Background()
	client, err := xai.NewClient(ctx)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	defer client.Close()
	stream, err := client.Responses.CreateStream(ctx, &xaiapiv1.GetCompletionsRequest{
		Model:    "grok-4.3",
		Messages: []*xaiapiv1.Message{messages.UserText("Say hi from integration test")},
	})
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	acc := responses.NewAccumulator()
	start := time.Now()
	var firstChunk time.Duration
	it := stream.Iterator(ctx)
	for {
		chunk, ok, err := it.Next()
		if errors.Is(err, context.Canceled) {
			t.Fatalf("context canceled")
		}
		// Surface real errors before ending on !ok: Next returns ok=false for
		// both a clean EOF and a mid-stream gRPC error.
		if err != nil && !errors.Is(err, io.EOF) {
			t.Fatalf("stream err: %v", err)
		}
		if !ok {
			break
		}
		if firstChunk == 0 {
			firstChunk = time.Since(start)
		}
		acc.AddChunk(chunk)
	}
	final := acc.Response()
	if len(final.GetOutputs()) == 0 {
		t.Fatalf("no outputs")
	}
	if final.GetUsage() == nil {
		t.Fatalf("usage missing")
	}
	tokenCount := final.GetUsage().GetCompletionTokens()
	totalDuration := time.Since(start)
	metrics := fmt.Sprintf("TTFT=%s total=%s tokens=%d", firstChunk, totalDuration, tokenCount)
	if tokenCount > 0 {
		tokPerSec := float64(tokenCount) / totalDuration.Seconds()
		metrics += fmt.Sprintf(" tok/s=%.2f", tokPerSec)
	}
	fmt.Println(metrics)
}
