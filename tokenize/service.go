package tokenize

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/internal/raw"
)

// Service exposes tokenization RPCs.
type Service struct {
	raw *raw.TokenizeClient
}

// NewService constructs a tokenize service.
func NewService(rawClient *raw.TokenizeClient) Service {
	return Service{raw: rawClient}
}

// TokenizeText converts text to tokens using the requested model.
func (s Service) TokenizeText(ctx context.Context, req *xaiapiv1.TokenizeTextRequest, opts ...grpc.CallOption) (*xaiapiv1.TokenizeTextResponse, error) {
	return s.raw.TokenizeText(ctx, req, opts...)
}
