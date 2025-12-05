package responses

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/modelrelay/xai-go/chat"
	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// Service exposes a gRPC-first Responses API layered on top of the Chat service.
type Service struct {
	chat chat.Service
}

// NewService builds a responses service from the shared Chat service.
func NewService(chatService chat.Service) Service {
	return Service{chat: chatService}
}

// Create issues a unary completion request and blocks until the full response is available.
func (s Service) Create(ctx context.Context, req *xaiapiv1.GetCompletionsRequest, opts ...grpc.CallOption) (*xaiapiv1.GetChatCompletionResponse, error) {
	return s.chat.GetCompletion(ctx, req, opts...)
}

// CreateStream issues a streaming completion request returning chunks as they are produced.
func (s Service) CreateStream(ctx context.Context, req *xaiapiv1.GetCompletionsRequest, opts ...grpc.CallOption) (*Stream, error) {
	stream, err := s.chat.GetCompletionChunk(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &Stream{raw: stream}, nil
}

// StartDeferred kicks off a deferred completion and immediately returns its request ID.
func (s Service) StartDeferred(ctx context.Context, req *xaiapiv1.GetCompletionsRequest, opts ...grpc.CallOption) (*xaiapiv1.StartDeferredResponse, error) {
	return s.chat.StartDeferredCompletion(ctx, req, opts...)
}

// GetDeferred fetches the result of a deferred completion by request ID.
func (s Service) GetDeferred(ctx context.Context, requestID string, opts ...grpc.CallOption) (*xaiapiv1.GetDeferredCompletionResponse, error) {
	if requestID == "" {
		return nil, fmt.Errorf("requestID is required")
	}
	return s.chat.GetDeferredCompletion(ctx, &xaiapiv1.GetDeferredRequest{RequestId: requestID}, opts...)
}

// PollDeferredCompletion polls a deferred request until it completes or the context is canceled.
func (s Service) PollDeferredCompletion(ctx context.Context, requestID string, interval time.Duration, opts ...grpc.CallOption) (*xaiapiv1.GetChatCompletionResponse, error) {
	if interval <= 0 {
		interval = 500 * time.Millisecond
	}
	resp, err := pollDeferred(ctx, interval, func(ctx context.Context) (*xaiapiv1.GetDeferredCompletionResponse, error) {
		return s.GetDeferred(ctx, requestID, opts...)
	})
	if err != nil {
		return nil, err
	}
	if resp.GetResponse() == nil {
		return nil, ErrDeferredNoResponse
	}
	return resp.GetResponse(), nil
}

// Retrieve fetches a stored completion by response ID.
func (s Service) Retrieve(ctx context.Context, responseID string, opts ...grpc.CallOption) (*xaiapiv1.GetChatCompletionResponse, error) {
	if responseID == "" {
		return nil, fmt.Errorf("responseID is required")
	}
	req := &xaiapiv1.GetStoredCompletionRequest{ResponseId: responseID}
	return s.chat.GetStoredCompletion(ctx, req, opts...)
}

// Delete removes a stored completion by response ID.
func (s Service) Delete(ctx context.Context, responseID string, opts ...grpc.CallOption) (*xaiapiv1.DeleteStoredCompletionResponse, error) {
	if responseID == "" {
		return nil, fmt.Errorf("responseID is required")
	}
	req := &xaiapiv1.DeleteStoredCompletionRequest{ResponseId: responseID}
	return s.chat.DeleteStoredCompletion(ctx, req, opts...)
}
