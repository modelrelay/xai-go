// Package batch provides access to xAI's batch management API.
package batch

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/internal/raw"
)

// Service exposes batch management RPCs.
type Service struct {
	raw *raw.BatchClient
}

// NewService constructs a batch service.
func NewService(rawClient *raw.BatchClient) Service {
	return Service{raw: rawClient}
}

func (s Service) CreateBatch(ctx context.Context, req *xaiapiv1.CreateBatchRequest, opts ...grpc.CallOption) (*xaiapiv1.Batch, error) {
	return s.raw.CreateBatch(ctx, req, opts...)
}

func (s Service) GetBatch(ctx context.Context, req *xaiapiv1.GetBatchRequest, opts ...grpc.CallOption) (*xaiapiv1.Batch, error) {
	return s.raw.GetBatch(ctx, req, opts...)
}

func (s Service) ListBatches(ctx context.Context, req *xaiapiv1.ListBatchesRequest, opts ...grpc.CallOption) (*xaiapiv1.ListBatchesResponse, error) {
	return s.raw.ListBatches(ctx, req, opts...)
}

func (s Service) CancelBatch(ctx context.Context, req *xaiapiv1.CancelBatchRequest, opts ...grpc.CallOption) (*xaiapiv1.Batch, error) {
	return s.raw.CancelBatch(ctx, req, opts...)
}

func (s Service) AddBatchRequests(ctx context.Context, req *xaiapiv1.AddBatchRequestsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return s.raw.AddBatchRequests(ctx, req, opts...)
}

func (s Service) ListBatchRequestMetadata(ctx context.Context, req *xaiapiv1.ListBatchRequestMetadataRequest, opts ...grpc.CallOption) (*xaiapiv1.ListBatchRequestMetadataResponse, error) {
	return s.raw.ListBatchRequestMetadata(ctx, req, opts...)
}

func (s Service) ListBatchResults(ctx context.Context, req *xaiapiv1.ListBatchResultsRequest, opts ...grpc.CallOption) (*xaiapiv1.ListBatchResultsResponse, error) {
	return s.raw.ListBatchResults(ctx, req, opts...)
}

func (s Service) GetBatchRequestResult(ctx context.Context, req *xaiapiv1.GetBatchRequestResultRequest, opts ...grpc.CallOption) (*xaiapiv1.GetBatchRequestResultResponse, error) {
	return s.raw.GetBatchRequestResult(ctx, req, opts...)
}
