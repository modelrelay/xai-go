package auth

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/internal/raw"
)

// Service exposes auth RPCs.
type Service struct {
	raw *raw.AuthClient
}

// NewService constructs an auth service.
func NewService(rawClient *raw.AuthClient) Service {
	return Service{raw: rawClient}
}

// GetAPIKeyInfo retrieves metadata for the caller's API key.
func (s Service) GetAPIKeyInfo(ctx context.Context, opts ...grpc.CallOption) (*xaiapiv1.ApiKey, error) {
	return s.raw.GetAPIKeyInfo(ctx, opts...)
}
