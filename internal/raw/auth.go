package raw

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// AuthClient wraps the auth service.
type AuthClient struct {
	stub xaiapiv1.AuthClient
}

// NewAuthClient creates a new auth client.
func NewAuthClient(stub xaiapiv1.AuthClient) *AuthClient {
	return &AuthClient{stub: stub}
}

// GetAPIKeyInfo fetches metadata about the current API key.
func (c *AuthClient) GetAPIKeyInfo(ctx context.Context, opts ...grpc.CallOption) (*xaiapiv1.ApiKey, error) {
	return c.stub.GetApiKeyInfo(ctx, &emptypb.Empty{}, opts...)
}
