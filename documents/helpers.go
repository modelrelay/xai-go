package documents

import (
	"encoding/json"
	"errors"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// CollectionSource builds a DocumentsSource scoped to the given collection IDs.
func CollectionSource(collectionIDs ...string) *xaiapiv1.DocumentsSource {
	return &xaiapiv1.DocumentsSource{
		CollectionIds: append([]string(nil), collectionIDs...),
	}
}

// SearchRequest constructs a SearchRequest with optional builders.
func SearchRequest(query string, source *xaiapiv1.DocumentsSource, opts ...SearchOption) *xaiapiv1.SearchRequest {
	req := &xaiapiv1.SearchRequest{
		Query:  query,
		Source: source,
	}
	for _, opt := range opts {
		opt(req)
	}
	return req
}

// SearchOption mutates a SearchRequest.
type SearchOption func(*xaiapiv1.SearchRequest)

// WithLimit sets the limit on results.
func WithLimit(limit int32) SearchOption {
	return func(req *xaiapiv1.SearchRequest) {
		req.Limit = &limit
	}
}

// WithRankingMetric sets the ranking metric to use.
func WithRankingMetric(metric xaiapiv1.RankingMetric) SearchOption {
	return func(req *xaiapiv1.SearchRequest) {
		req.RankingMetric = &metric
	}
}

// WithInstructions sets optional instructions.
func WithInstructions(text string) SearchOption {
	return func(req *xaiapiv1.SearchRequest) {
		req.Instructions = &text
	}
}

// ToolMatch summarizes a search match for tool responses.
type ToolMatch struct {
	FileID        string   `json:"file_id"`
	ChunkID       string   `json:"chunk_id"`
	ChunkContent  string   `json:"chunk_content"`
	Score         float32  `json:"score"`
	CollectionIDs []string `json:"collection_ids"`
}

type toolPayload struct {
	ToolCallID string      `json:"tool_call_id,omitempty"`
	Matches    []ToolMatch `json:"matches"`
}

// ToolResponsePayload serializes matches into a JSON payload suitable for ROLE_TOOL messages.
func ToolResponsePayload(toolCallID string, matches []*xaiapiv1.SearchMatch, limit int) (string, error) {
	if len(matches) == 0 {
		return "", errors.New("matches cannot be empty")
	}
	if limit <= 0 || limit > len(matches) {
		limit = len(matches)
	}
	payload := toolPayload{ToolCallID: toolCallID}
	for i := 0; i < limit; i++ {
		m := matches[i]
		payload.Matches = append(payload.Matches, ToolMatch{
			FileID:        m.GetFileId(),
			ChunkID:       m.GetChunkId(),
			ChunkContent:  m.GetChunkContent(),
			Score:         m.GetScore(),
			CollectionIDs: append([]string(nil), m.GetCollectionIds()...),
		})
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToolResponseMessage creates a ROLE_TOOL message containing search matches in JSON form.
func ToolResponseMessage(toolCallID string, matches []*xaiapiv1.SearchMatch, limit int) (*xaiapiv1.Message, error) {
	payload, err := ToolResponsePayload(toolCallID, matches, limit)
	if err != nil {
		return nil, err
	}
	return &xaiapiv1.Message{
		Role: xaiapiv1.MessageRole_ROLE_TOOL,
		Content: []*xaiapiv1.Content{
			{Content: &xaiapiv1.Content_Text{Text: payload}},
		},
	}, nil
}
