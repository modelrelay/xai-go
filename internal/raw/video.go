package raw

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// VideoClient wraps the generated video service.
type VideoClient struct {
	stub xaiapiv1.VideoClient
}

// NewVideoClient creates a new raw video client.
func NewVideoClient(stub xaiapiv1.VideoClient) *VideoClient {
	return &VideoClient{stub: stub}
}

func (c *VideoClient) GenerateVideo(ctx context.Context, req *xaiapiv1.GenerateVideoRequest, opts ...grpc.CallOption) (*xaiapiv1.StartDeferredResponse, error) {
	return c.stub.GenerateVideo(ctx, req, opts...)
}

func (c *VideoClient) ExtendVideo(ctx context.Context, req *xaiapiv1.ExtendVideoRequest, opts ...grpc.CallOption) (*xaiapiv1.StartDeferredResponse, error) {
	return c.stub.ExtendVideo(ctx, req, opts...)
}

func (c *VideoClient) GetDeferredVideo(ctx context.Context, req *xaiapiv1.GetDeferredVideoRequest, opts ...grpc.CallOption) (*xaiapiv1.GetDeferredVideoResponse, error) {
	return c.stub.GetDeferredVideo(ctx, req, opts...)
}
