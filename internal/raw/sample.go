package raw

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// SampleClient wraps the Sample service stub.
type SampleClient struct {
	stub     xaiapiv1.SampleClient
	defaults Defaults
}

// NewSampleClient creates a raw Sample client.
func NewSampleClient(stub xaiapiv1.SampleClient, defaults Defaults) *SampleClient {
	return &SampleClient{stub: stub, defaults: defaults}
}

// SampleText invokes the unary sampling endpoint.
func (c *SampleClient) SampleText(ctx context.Context, req *xaiapiv1.SampleTextRequest, opts ...grpc.CallOption) (*xaiapiv1.SampleTextResponse, error) {
	c.applyDefaults(req)
	return c.stub.SampleText(ctx, req, opts...)
}

// SampleTextStreaming opens a streaming sampler.
func (c *SampleClient) SampleTextStreaming(ctx context.Context, req *xaiapiv1.SampleTextRequest, opts ...grpc.CallOption) (xaiapiv1.Sample_SampleTextStreamingClient, error) {
	c.applyDefaults(req)
	return c.stub.SampleTextStreaming(ctx, req, opts...)
}

func (c *SampleClient) applyDefaults(req *xaiapiv1.SampleTextRequest) {
	if req == nil {
		return
	}
	if req.User == "" && c.defaults.User != "" {
		req.User = c.defaults.User
	}
}
