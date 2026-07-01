package documents

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/internal/raw"
)

// Service exposes the documents search RPC.
type Service struct {
	raw *raw.DocumentsClient
}

// NewService constructs a documents service.
func NewService(rawClient *raw.DocumentsClient) Service {
	return Service{raw: rawClient}
}

// Search performs a document search.
func (s Service) Search(ctx context.Context, req *xaiapiv1.SearchRequest, opts ...grpc.CallOption) (*xaiapiv1.SearchResponse, error) {
	return s.raw.Search(ctx, req, opts...)
}
