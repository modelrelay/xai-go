package embeddings

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/internal/raw"
)

// Service exposes the Embedder RPCs.
type Service struct {
	raw *raw.EmbeddingClient
}

// NewService constructs an embeddings service.
func NewService(rawClient *raw.EmbeddingClient) Service {
	return Service{raw: rawClient}
}

// Embed produces embeddings for the provided inputs.
func (s Service) Embed(ctx context.Context, req *xaiapiv1.EmbedRequest, opts ...grpc.CallOption) (*xaiapiv1.EmbedResponse, error) {
	return s.raw.Embed(ctx, req, opts...)
}
