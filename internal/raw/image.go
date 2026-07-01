package raw

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// ImageClient wraps the image generation service.
type ImageClient struct {
	stub xaiapiv1.ImageClient
}

// NewImageClient creates a new image client.
func NewImageClient(stub xaiapiv1.ImageClient) *ImageClient {
	return &ImageClient{stub: stub}
}

// GenerateImage calls the image generation RPC.
func (c *ImageClient) GenerateImage(ctx context.Context, req *xaiapiv1.GenerateImageRequest, opts ...grpc.CallOption) (*xaiapiv1.ImageResponse, error) {
	return c.stub.GenerateImage(ctx, req, opts...)
}
