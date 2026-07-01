package images

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/internal/raw"
)

// Service exposes image generation RPCs.
type Service struct {
	raw *raw.ImageClient
}

// NewService constructs an image service.
func NewService(rawClient *raw.ImageClient) Service {
	return Service{raw: rawClient}
}

// GenerateImage creates images from prompts/image inputs.
func (s Service) GenerateImage(ctx context.Context, req *xaiapiv1.GenerateImageRequest, opts ...grpc.CallOption) (*xaiapiv1.ImageResponse, error) {
	return s.raw.GenerateImage(ctx, req, opts...)
}
