package raw

import (
	"context"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// FilesClient wraps the generated files service.
type FilesClient struct {
	stub xaiapiv1.FilesClient
}

// NewFilesClient creates a new raw files client.
func NewFilesClient(stub xaiapiv1.FilesClient) *FilesClient {
	return &FilesClient{stub: stub}
}

// UploadFile opens a client stream for uploading a file in chunks.
func (c *FilesClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[xaiapiv1.UploadFileChunk, xaiapiv1.File], error) {
	return c.stub.UploadFile(ctx, opts...)
}

func (c *FilesClient) ListFiles(ctx context.Context, req *xaiapiv1.ListFilesRequest, opts ...grpc.CallOption) (*xaiapiv1.ListFilesResponse, error) {
	return c.stub.ListFiles(ctx, req, opts...)
}

func (c *FilesClient) RetrieveFile(ctx context.Context, req *xaiapiv1.RetrieveFileRequest, opts ...grpc.CallOption) (*xaiapiv1.File, error) {
	return c.stub.RetrieveFile(ctx, req, opts...)
}

func (c *FilesClient) DeleteFile(ctx context.Context, req *xaiapiv1.DeleteFileRequest, opts ...grpc.CallOption) (*xaiapiv1.DeleteFileResponse, error) {
	return c.stub.DeleteFile(ctx, req, opts...)
}

// RetrieveFileContent opens a server stream of file content chunks.
func (c *FilesClient) RetrieveFileContent(ctx context.Context, req *xaiapiv1.RetrieveFileContentRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[xaiapiv1.FileContentChunk], error) {
	return c.stub.RetrieveFileContent(ctx, req, opts...)
}
