package raw

import (
	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	billingv1 "github.com/modelrelay/xai-go/gen/xai/management_api/v1"
)

// Clients aggregates typed raw clients for every proto service.
type Clients struct {
	Chat    *ChatClient
	Sample  *SampleClient
	Embed   *EmbeddingClient
	Image   *ImageClient
	Docs    *DocumentsClient
	Token   *TokenizeClient
	Models  *ModelsClient
	Auth    *AuthClient
	Billing *BillingClient
}

// NewClients constructs raw clients powered by the provided connection.
func NewClients(conn *grpc.ClientConn, defaults Defaults) Clients {
	return Clients{
		Chat:    NewChatClient(xaiapiv1.NewChatClient(conn), defaults),
		Sample:  NewSampleClient(xaiapiv1.NewSampleClient(conn), defaults),
		Embed:   NewEmbeddingClient(xaiapiv1.NewEmbedderClient(conn)),
		Image:   NewImageClient(xaiapiv1.NewImageClient(conn)),
		Docs:    NewDocumentsClient(xaiapiv1.NewDocumentsClient(conn)),
		Token:   NewTokenizeClient(xaiapiv1.NewTokenizeClient(conn), defaults),
		Models:  NewModelsClient(xaiapiv1.NewModelsClient(conn)),
		Auth:    NewAuthClient(xaiapiv1.NewAuthClient(conn)),
		Billing: NewBillingClient(billingv1.NewUISvcClient(conn)),
	}
}
