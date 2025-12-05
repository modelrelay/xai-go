package responses

import (
	"testing"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

func TestRequireStoredRequests(t *testing.T) {
	req := &xaiapiv1.GetCompletionsRequest{}
	if err := RequireStoredRequests(req, "resp_1"); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !req.GetStoreMessages() {
		t.Fatalf("store_messages not set")
	}
	if req.GetPreviousResponseId() != "resp_1" {
		t.Fatalf("previous id mismatch")
	}
}
