package messages

import (
	"testing"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

func TestUserText(t *testing.T) {
	msg := UserText("hello")
	if msg.GetRole() != xaiapiv1.MessageRole_ROLE_USER {
		t.Fatalf("role mismatch: %v", msg.GetRole())
	}
	if got := msg.GetContent()[0].GetText(); got != "hello" {
		t.Fatalf("content mismatch: %s", got)
	}
}

func TestImageURL(t *testing.T) {
	msg := AssistantText("see image")
	ImageURL(msg, "https://example.com/cat.png", xaiapiv1.ImageDetail_DETAIL_HIGH)
	if len(msg.GetContent()) != 2 {
		t.Fatalf("content len mismatch: %d", len(msg.GetContent()))
	}
	img := msg.GetContent()[1].GetImageUrl()
	if img.GetDetail() != xaiapiv1.ImageDetail_DETAIL_HIGH {
		t.Fatalf("detail mismatch: %v", img.GetDetail())
	}
}
