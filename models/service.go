package models

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/internal/raw"
)

// Service exposes model discovery RPCs.
type Service struct {
	raw *raw.ModelsClient
}

// NewService constructs a models service.
func NewService(rawClient *raw.ModelsClient) Service {
	return Service{raw: rawClient}
}

func (s Service) ListLanguageModels(ctx context.Context, opts ...grpc.CallOption) (*xaiapiv1.ListLanguageModelsResponse, error) {
	return s.raw.ListLanguageModels(ctx, opts...)
}

func (s Service) ListEmbeddingModels(ctx context.Context, opts ...grpc.CallOption) (*xaiapiv1.ListEmbeddingModelsResponse, error) {
	return s.raw.ListEmbeddingModels(ctx, opts...)
}

func (s Service) ListImageGenerationModels(ctx context.Context, opts ...grpc.CallOption) (*xaiapiv1.ListImageGenerationModelsResponse, error) {
	return s.raw.ListImageGenerationModels(ctx, opts...)
}

func (s Service) GetLanguageModel(ctx context.Context, req *xaiapiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiapiv1.LanguageModel, error) {
	return s.raw.GetLanguageModel(ctx, req, opts...)
}

func (s Service) GetEmbeddingModel(ctx context.Context, req *xaiapiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiapiv1.EmbeddingModel, error) {
	return s.raw.GetEmbeddingModel(ctx, req, opts...)
}

func (s Service) GetImageGenerationModel(ctx context.Context, req *xaiapiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiapiv1.ImageGenerationModel, error) {
	return s.raw.GetImageGenerationModel(ctx, req, opts...)
}
