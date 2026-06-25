package raw

import (
	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// Clients aggregates typed raw clients for every proto service.
type Clients struct {
	Chat   *ChatClient
	Sample *SampleClient
	Embed  *EmbeddingClient
	Image  *ImageClient
	Docs   *DocumentsClient
	Token  *TokenizeClient
	Models *ModelsClient
	Auth   *AuthClient
	Batch  *BatchClient
	Files  *FilesClient
	Video  *VideoClient
}

// NewClients constructs raw clients powered by the provided connection.
func NewClients(conn *grpc.ClientConn, defaults Defaults) Clients {
	return Clients{
		Chat:   NewChatClient(xaiapiv1.NewChatClient(conn), defaults),
		Sample: NewSampleClient(xaiapiv1.NewSampleClient(conn), defaults),
		Embed:  NewEmbeddingClient(xaiapiv1.NewEmbedderClient(conn)),
		Image:  NewImageClient(xaiapiv1.NewImageClient(conn)),
		Docs:   NewDocumentsClient(xaiapiv1.NewDocumentsClient(conn)),
		Token:  NewTokenizeClient(xaiapiv1.NewTokenizeClient(conn), defaults),
		Models: NewModelsClient(xaiapiv1.NewModelsClient(conn)),
		Auth:   NewAuthClient(xaiapiv1.NewAuthClient(conn)),
		Batch:  NewBatchClient(xaiapiv1.NewBatchMgmtClient(conn)),
		Files:  NewFilesClient(xaiapiv1.NewFilesClient(conn)),
		Video:  NewVideoClient(xaiapiv1.NewVideoClient(conn)),
	}
}
