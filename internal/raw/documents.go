package raw

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// DocumentsClient wraps the documents search service.
type DocumentsClient struct {
	stub xaiapiv1.DocumentsClient
}

// NewDocumentsClient creates a new documents client.
func NewDocumentsClient(stub xaiapiv1.DocumentsClient) *DocumentsClient {
	return &DocumentsClient{stub: stub}
}

// Search calls the document search RPC.
func (c *DocumentsClient) Search(ctx context.Context, req *xaiapiv1.SearchRequest, opts ...grpc.CallOption) (*xaiapiv1.SearchResponse, error) {
	return c.stub.Search(ctx, req, opts...)
}
