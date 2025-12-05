package encrypted

import (
	"errors"
	"testing"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

func TestDecryptResponse(t *testing.T) {
	resp := &xaiapiv1.GetChatCompletionResponse{
		Outputs: []*xaiapiv1.CompletionOutput{
			{Message: &xaiapiv1.CompletionMessage{EncryptedContent: "secret"}},
		},
	}
	decode := func(s string) (string, error) { return "plain:" + s, nil }
	if err := DecryptResponse(resp, decode); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.GetOutputs()[0].GetMessage().GetContent() != "plain:secret" {
		t.Fatalf("content mismatch: %s", resp.GetOutputs()[0].GetMessage().GetContent())
	}
}

func TestDecryptChunkError(t *testing.T) {
	chunk := &xaiapiv1.GetChatCompletionChunk{
		Outputs: []*xaiapiv1.CompletionOutputChunk{
			{Delta: &xaiapiv1.Delta{EncryptedContent: "fail"}},
		},
	}
	decode := func(string) (string, error) { return "", errors.New("boom") }
	if err := DecryptChunk(chunk, decode); err == nil {
		t.Fatalf("expected error")
	}
}
