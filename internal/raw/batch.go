package raw

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// BatchClient wraps the generated batch management service.
type BatchClient struct {
	stub xaiapiv1.BatchMgmtClient
}

// NewBatchClient creates a new raw batch client.
func NewBatchClient(stub xaiapiv1.BatchMgmtClient) *BatchClient {
	return &BatchClient{stub: stub}
}

func (c *BatchClient) CreateBatch(ctx context.Context, req *xaiapiv1.CreateBatchRequest, opts ...grpc.CallOption) (*xaiapiv1.Batch, error) {
	return c.stub.CreateBatch(ctx, req, opts...)
}

func (c *BatchClient) GetBatch(ctx context.Context, req *xaiapiv1.GetBatchRequest, opts ...grpc.CallOption) (*xaiapiv1.Batch, error) {
	return c.stub.GetBatch(ctx, req, opts...)
}

func (c *BatchClient) ListBatches(ctx context.Context, req *xaiapiv1.ListBatchesRequest, opts ...grpc.CallOption) (*xaiapiv1.ListBatchesResponse, error) {
	return c.stub.ListBatches(ctx, req, opts...)
}

func (c *BatchClient) CancelBatch(ctx context.Context, req *xaiapiv1.CancelBatchRequest, opts ...grpc.CallOption) (*xaiapiv1.Batch, error) {
	return c.stub.CancelBatch(ctx, req, opts...)
}

func (c *BatchClient) AddBatchRequests(ctx context.Context, req *xaiapiv1.AddBatchRequestsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.stub.AddBatchRequests(ctx, req, opts...)
}

func (c *BatchClient) ListBatchRequestMetadata(ctx context.Context, req *xaiapiv1.ListBatchRequestMetadataRequest, opts ...grpc.CallOption) (*xaiapiv1.ListBatchRequestMetadataResponse, error) {
	return c.stub.ListBatchRequestMetadata(ctx, req, opts...)
}

func (c *BatchClient) ListBatchResults(ctx context.Context, req *xaiapiv1.ListBatchResultsRequest, opts ...grpc.CallOption) (*xaiapiv1.ListBatchResultsResponse, error) {
	return c.stub.ListBatchResults(ctx, req, opts...)
}

func (c *BatchClient) GetBatchRequestResult(ctx context.Context, req *xaiapiv1.GetBatchRequestResultRequest, opts ...grpc.CallOption) (*xaiapiv1.GetBatchRequestResultResponse, error) {
	return c.stub.GetBatchRequestResult(ctx, req, opts...)
}
