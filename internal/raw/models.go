package raw

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ModelsClient wraps the models service.
type ModelsClient struct {
	stub xaiapiv1.ModelsClient
}

// NewModelsClient creates a new models client.
func NewModelsClient(stub xaiapiv1.ModelsClient) *ModelsClient {
	return &ModelsClient{stub: stub}
}

func (c *ModelsClient) ListLanguageModels(ctx context.Context, opts ...grpc.CallOption) (*xaiapiv1.ListLanguageModelsResponse, error) {
	return c.stub.ListLanguageModels(ctx, &emptypb.Empty{}, opts...)
}

func (c *ModelsClient) ListEmbeddingModels(ctx context.Context, opts ...grpc.CallOption) (*xaiapiv1.ListEmbeddingModelsResponse, error) {
	return c.stub.ListEmbeddingModels(ctx, &emptypb.Empty{}, opts...)
}

func (c *ModelsClient) ListImageGenerationModels(ctx context.Context, opts ...grpc.CallOption) (*xaiapiv1.ListImageGenerationModelsResponse, error) {
	return c.stub.ListImageGenerationModels(ctx, &emptypb.Empty{}, opts...)
}

func (c *ModelsClient) GetLanguageModel(ctx context.Context, req *xaiapiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiapiv1.LanguageModel, error) {
	return c.stub.GetLanguageModel(ctx, req, opts...)
}

func (c *ModelsClient) GetEmbeddingModel(ctx context.Context, req *xaiapiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiapiv1.EmbeddingModel, error) {
	return c.stub.GetEmbeddingModel(ctx, req, opts...)
}

func (c *ModelsClient) GetImageGenerationModel(ctx context.Context, req *xaiapiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiapiv1.ImageGenerationModel, error) {
	return c.stub.GetImageGenerationModel(ctx, req, opts...)
}
