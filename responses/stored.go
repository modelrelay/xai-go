package responses

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// RetrieveAndDelete fetches a stored response and deletes it afterwards.
func (s Service) RetrieveAndDelete(ctx context.Context, responseID string, opts ...grpc.CallOption) (*xaiapiv1.GetChatCompletionResponse, error) {
	resp, err := s.Retrieve(ctx, responseID, opts...)
	if err != nil {
		return nil, err
	}
	if _, err := s.Delete(ctx, responseID, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// RequireStoredRequests sets store_messages=true and validates response IDs when chaining conversations.
func RequireStoredRequests(req *xaiapiv1.GetCompletionsRequest, previousResponseID string) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	req.StoreMessages = true
	if previousResponseID != "" {
		req.PreviousResponseId = &previousResponseID
	}
	return nil
}
