package chat

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/internal/raw"
)

// Service exposes high-level helpers for the Chat API.
type Service struct {
	raw *raw.ChatClient
}

// NewService wires a Chat service backed by the raw client.
func NewService(rawClient *raw.ChatClient) Service {
	return Service{raw: rawClient}
}

// GetCompletion requests a fully-buffered completion response.
func (s Service) GetCompletion(ctx context.Context, req *xaiapiv1.GetCompletionsRequest, opts ...grpc.CallOption) (*xaiapiv1.GetChatCompletionResponse, error) {
	return s.raw.GetCompletion(ctx, req, opts...)
}

// GetCompletionChunk opens a streaming completion session.
func (s Service) GetCompletionChunk(ctx context.Context, req *xaiapiv1.GetCompletionsRequest, opts ...grpc.CallOption) (xaiapiv1.Chat_GetCompletionChunkClient, error) {
	return s.raw.GetCompletionChunk(ctx, req, opts...)
}

// StartDeferredCompletion starts an async completion and returns the request ID.
func (s Service) StartDeferredCompletion(ctx context.Context, req *xaiapiv1.GetCompletionsRequest, opts ...grpc.CallOption) (*xaiapiv1.StartDeferredResponse, error) {
	return s.raw.StartDeferredCompletion(ctx, req, opts...)
}

// GetDeferredCompletion fetches the result of a deferred completion.
func (s Service) GetDeferredCompletion(ctx context.Context, req *xaiapiv1.GetDeferredRequest, opts ...grpc.CallOption) (*xaiapiv1.GetDeferredCompletionResponse, error) {
	return s.raw.GetDeferredCompletion(ctx, req, opts...)
}

// GetStoredCompletion retrieves a stored completion by ID.
func (s Service) GetStoredCompletion(ctx context.Context, req *xaiapiv1.GetStoredCompletionRequest, opts ...grpc.CallOption) (*xaiapiv1.GetChatCompletionResponse, error) {
	return s.raw.GetStoredCompletion(ctx, req, opts...)
}

// DeleteStoredCompletion deletes a stored completion by ID.
func (s Service) DeleteStoredCompletion(ctx context.Context, req *xaiapiv1.DeleteStoredCompletionRequest, opts ...grpc.CallOption) (*xaiapiv1.DeleteStoredCompletionResponse, error) {
	return s.raw.DeleteStoredCompletion(ctx, req, opts...)
}

// CompactContext compacts a conversation's context.
func (s Service) CompactContext(ctx context.Context, req *xaiapiv1.CompactContextRequest, opts ...grpc.CallOption) (*xaiapiv1.CompactContextResponse, error) {
	return s.raw.CompactContext(ctx, req, opts...)
}
