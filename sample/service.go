package sample

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/internal/raw"
)

// Service exposes the low-level Sample API.
type Service struct {
	raw *raw.SampleClient
}

// NewService builds a Sample service backed by raw client.
func NewService(rawClient *raw.SampleClient) Service {
	return Service{raw: rawClient}
}

// SampleText performs unary sampling.
func (s Service) SampleText(ctx context.Context, req *xaiapiv1.SampleTextRequest, opts ...grpc.CallOption) (*xaiapiv1.SampleTextResponse, error) {
	return s.raw.SampleText(ctx, req, opts...)
}

// SampleTextStreaming streams sampling results chunk by chunk.
func (s Service) SampleTextStreaming(ctx context.Context, req *xaiapiv1.SampleTextRequest, opts ...grpc.CallOption) (xaiapiv1.Sample_SampleTextStreamingClient, error) {
	return s.raw.SampleTextStreaming(ctx, req, opts...)
}
