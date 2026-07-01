package raw

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// EmbeddingClient wraps the Embedder service.
type EmbeddingClient struct {
	stub xaiapiv1.EmbedderClient
}

// NewEmbeddingClient creates a new embedding client.
func NewEmbeddingClient(stub xaiapiv1.EmbedderClient) *EmbeddingClient {
	return &EmbeddingClient{stub: stub}
}

// Embed invokes the embedding RPC.
func (c *EmbeddingClient) Embed(ctx context.Context, req *xaiapiv1.EmbedRequest, opts ...grpc.CallOption) (*xaiapiv1.EmbedResponse, error) {
	return c.stub.Embed(ctx, req, opts...)
}
