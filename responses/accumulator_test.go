package responses

import (
	"testing"
	"time"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAccumulatorMergesChunks(t *testing.T) {
	acc := NewAccumulator()

	chunk1 := &xaiapiv1.GetChatCompletionChunk{
		Id:      "resp_123",
		Model:   "grok-2",
		Created: timestamppb.New(time.Unix(42, 0)),
		Outputs: []*xaiapiv1.CompletionOutputChunk{
			{
				Index: 0,
				Delta: &xaiapiv1.Delta{
					Content: "Hello",
				},
				Logprobs: &xaiapiv1.LogProbs{
					Content: []*xaiapiv1.LogProb{
						{Token: "Hello"},
					},
				},
			},
		},
	}
	acc.AddChunk(chunk1)

	chunk2 := &xaiapiv1.GetChatCompletionChunk{
		Id:    "resp_123",
		Model: "grok-2",
		Outputs: []*xaiapiv1.CompletionOutputChunk{
			{
				Index:        0,
				Delta:        &xaiapiv1.Delta{Content: " world"},
				FinishReason: xaiapiv1.FinishReason_REASON_STOP,
			},
		},
		Citations: []string{"https://example.com"},
	}
	acc.AddChunk(chunk2)

	resp := acc.Response()

	if got, want := resp.GetId(), "resp_123"; got != want {
		t.Fatalf("id mismatch: got %s want %s", got, want)
	}
	if got, want := len(resp.GetOutputs()), 1; got != want {
		t.Fatalf("outputs len mismatch: got %d want %d", got, want)
	}
	output := resp.GetOutputs()[0]
	if got, want := output.GetMessage().GetContent(), "Hello world"; got != want {
		t.Fatalf("content mismatch: got %q want %q", got, want)
	}
	if got := output.GetFinishReason(); got != xaiapiv1.FinishReason_REASON_STOP {
		t.Fatalf("finish reason mismatch: got %v", got)
	}
	if got, want := len(resp.GetCitations()), 1; got != want {
		t.Fatalf("citations len mismatch: got %d want %d", got, want)
	}
	if got := resp.GetCitations()[0]; got != "https://example.com" {
		t.Fatalf("citation mismatch: got %s", got)
	}
	if got, want := len(output.GetLogprobs().GetContent()), 1; got != want {
		t.Fatalf("logprobs mismatch: got %d want %d", got, want)
	}
}

func BenchmarkAccumulatorAddChunk(b *testing.B) {
	chunk := &xaiapiv1.GetChatCompletionChunk{
		Id: "resp",
		Outputs: []*xaiapiv1.CompletionOutputChunk{
			{Index: 0, Delta: &xaiapiv1.Delta{Content: "hello world"}},
		},
	}
	for i := 0; i < b.N; i++ {
		acc := NewAccumulator()
		acc.AddChunk(chunk)
		_ = acc.Response()
	}
}
