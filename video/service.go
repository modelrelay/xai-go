// Package video provides access to xAI's video generation API.
package video

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/internal/raw"
)

// Service exposes video generation RPCs.
type Service struct {
	raw *raw.VideoClient
}

// NewService constructs a video service.
func NewService(rawClient *raw.VideoClient) Service {
	return Service{raw: rawClient}
}

func (s Service) GenerateVideo(ctx context.Context, req *xaiapiv1.GenerateVideoRequest, opts ...grpc.CallOption) (*xaiapiv1.StartDeferredResponse, error) {
	return s.raw.GenerateVideo(ctx, req, opts...)
}

func (s Service) ExtendVideo(ctx context.Context, req *xaiapiv1.ExtendVideoRequest, opts ...grpc.CallOption) (*xaiapiv1.StartDeferredResponse, error) {
	return s.raw.ExtendVideo(ctx, req, opts...)
}

func (s Service) GetDeferredVideo(ctx context.Context, req *xaiapiv1.GetDeferredVideoRequest, opts ...grpc.CallOption) (*xaiapiv1.GetDeferredVideoResponse, error) {
	return s.raw.GetDeferredVideo(ctx, req, opts...)
}
