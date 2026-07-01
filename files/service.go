// Package files provides access to xAI's files API.
package files

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/internal/raw"
)

// Service exposes file management RPCs.
type Service struct {
	raw *raw.FilesClient
}

// NewService constructs a files service.
func NewService(rawClient *raw.FilesClient) Service {
	return Service{raw: rawClient}
}

// UploadFile opens a client stream for uploading a file in chunks.
func (s Service) UploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[xaiapiv1.UploadFileChunk, xaiapiv1.File], error) {
	return s.raw.UploadFile(ctx, opts...)
}

func (s Service) ListFiles(ctx context.Context, req *xaiapiv1.ListFilesRequest, opts ...grpc.CallOption) (*xaiapiv1.ListFilesResponse, error) {
	return s.raw.ListFiles(ctx, req, opts...)
}

func (s Service) RetrieveFile(ctx context.Context, req *xaiapiv1.RetrieveFileRequest, opts ...grpc.CallOption) (*xaiapiv1.File, error) {
	return s.raw.RetrieveFile(ctx, req, opts...)
}

func (s Service) DeleteFile(ctx context.Context, req *xaiapiv1.DeleteFileRequest, opts ...grpc.CallOption) (*xaiapiv1.DeleteFileResponse, error) {
	return s.raw.DeleteFile(ctx, req, opts...)
}

// RetrieveFileContent opens a server stream of file content chunks.
func (s Service) RetrieveFileContent(ctx context.Context, req *xaiapiv1.RetrieveFileContentRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[xaiapiv1.FileContentChunk], error) {
	return s.raw.RetrieveFileContent(ctx, req, opts...)
}
