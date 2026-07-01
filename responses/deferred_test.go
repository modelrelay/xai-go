package responses

import (
	"context"
	"errors"
	"testing"
	"time"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

func TestPollDeferredCompletes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	calls := 0
	fetch := func(context.Context) (*xaiapiv1.GetDeferredCompletionResponse, error) {
		calls++
		if calls < 3 {
			return &xaiapiv1.GetDeferredCompletionResponse{Status: xaiapiv1.DeferredStatus_PENDING}, nil
		}
		return &xaiapiv1.GetDeferredCompletionResponse{
			Status: xaiapiv1.DeferredStatus_DONE,
			Response: &xaiapiv1.GetChatCompletionResponse{
				Id: "resp_123",
			},
		}, nil
	}

	resp, err := pollDeferred(ctx, 10*time.Millisecond, fetch)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.GetResponse().GetId() != "resp_123" {
		t.Fatalf("id mismatch: %s", resp.GetResponse().GetId())
	}
}

func TestPollDeferredExpires(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	fetch := func(context.Context) (*xaiapiv1.GetDeferredCompletionResponse, error) {
		return &xaiapiv1.GetDeferredCompletionResponse{Status: xaiapiv1.DeferredStatus_EXPIRED}, nil
	}

	if _, err := pollDeferred(ctx, 5*time.Millisecond, fetch); !errors.Is(err, ErrDeferredExpired) {
		t.Fatalf("expected ErrDeferredExpired, got %v", err)
	}
}

func TestPollDeferredContextCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	fetch := func(context.Context) (*xaiapiv1.GetDeferredCompletionResponse, error) {
		return &xaiapiv1.GetDeferredCompletionResponse{Status: xaiapiv1.DeferredStatus_PENDING}, nil
	}

	_, err := pollDeferred(ctx, 20*time.Millisecond, fetch)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline, got %v", err)
	}
}
