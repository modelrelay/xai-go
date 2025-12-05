package raw

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// TokenizeClient wraps the tokenize service.
type TokenizeClient struct {
	stub     xaiapiv1.TokenizeClient
	defaults Defaults
}

// NewTokenizeClient creates a new tokenize client.
func NewTokenizeClient(stub xaiapiv1.TokenizeClient, defaults Defaults) *TokenizeClient {
	return &TokenizeClient{stub: stub, defaults: defaults}
}

// TokenizeText invokes the tokenization RPC.
func (c *TokenizeClient) TokenizeText(ctx context.Context, req *xaiapiv1.TokenizeTextRequest, opts ...grpc.CallOption) (*xaiapiv1.TokenizeTextResponse, error) {
	c.applyDefaults(req)
	return c.stub.TokenizeText(ctx, req, opts...)
}

func (c *TokenizeClient) applyDefaults(req *xaiapiv1.TokenizeTextRequest) {
	if req == nil {
		return
	}
	if req.User == "" && c.defaults.User != "" {
		req.User = c.defaults.User
	}
}
