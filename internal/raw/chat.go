package raw

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// ChatClient wraps the generated ChatClient to inject request defaults.
type ChatClient struct {
	stub     xaiapiv1.ChatClient
	defaults Defaults
}

// NewChatClient creates a new raw Chat client.
func NewChatClient(stub xaiapiv1.ChatClient, defaults Defaults) *ChatClient {
	return &ChatClient{stub: stub, defaults: defaults}
}

// GetCompletion calls the unary completion RPC.
func (c *ChatClient) GetCompletion(ctx context.Context, req *xaiapiv1.GetCompletionsRequest, opts ...grpc.CallOption) (*xaiapiv1.GetChatCompletionResponse, error) {
	c.applyDefaults(req)
	return c.stub.GetCompletion(ctx, req, opts...)
}

// GetCompletionChunk starts streaming completion chunks.
func (c *ChatClient) GetCompletionChunk(ctx context.Context, req *xaiapiv1.GetCompletionsRequest, opts ...grpc.CallOption) (xaiapiv1.Chat_GetCompletionChunkClient, error) {
	c.applyDefaults(req)
	return c.stub.GetCompletionChunk(ctx, req, opts...)
}

// StartDeferredCompletion kicks off an asynchronous completion request.
func (c *ChatClient) StartDeferredCompletion(ctx context.Context, req *xaiapiv1.GetCompletionsRequest, opts ...grpc.CallOption) (*xaiapiv1.StartDeferredResponse, error) {
	c.applyDefaults(req)
	return c.stub.StartDeferredCompletion(ctx, req, opts...)
}

// GetDeferredCompletion polls a deferred completion.
func (c *ChatClient) GetDeferredCompletion(ctx context.Context, req *xaiapiv1.GetDeferredRequest, opts ...grpc.CallOption) (*xaiapiv1.GetDeferredCompletionResponse, error) {
	return c.stub.GetDeferredCompletion(ctx, req, opts...)
}

// GetStoredCompletion fetches a stored response.
func (c *ChatClient) GetStoredCompletion(ctx context.Context, req *xaiapiv1.GetStoredCompletionRequest, opts ...grpc.CallOption) (*xaiapiv1.GetChatCompletionResponse, error) {
	return c.stub.GetStoredCompletion(ctx, req, opts...)
}

// DeleteStoredCompletion removes a stored response.
func (c *ChatClient) DeleteStoredCompletion(ctx context.Context, req *xaiapiv1.DeleteStoredCompletionRequest, opts ...grpc.CallOption) (*xaiapiv1.DeleteStoredCompletionResponse, error) {
	return c.stub.DeleteStoredCompletion(ctx, req, opts...)
}

// CompactContext compacts a conversation's context.
func (c *ChatClient) CompactContext(ctx context.Context, req *xaiapiv1.CompactContextRequest, opts ...grpc.CallOption) (*xaiapiv1.CompactContextResponse, error) {
	return c.stub.CompactContext(ctx, req, opts...)
}

func (c *ChatClient) applyDefaults(req *xaiapiv1.GetCompletionsRequest) {
	if req == nil {
		return
	}
	if req.User == "" && c.defaults.User != "" {
		req.User = c.defaults.User
	}
}
