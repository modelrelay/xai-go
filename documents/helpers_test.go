package documents

import (
	"testing"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

func TestSearchBuilders(t *testing.T) {
	source := CollectionSource("col1", "col2")
	req := SearchRequest("hello", source, WithLimit(5), WithRankingMetric(xaiapiv1.RankingMetric_RANKING_METRIC_COSINE_SIMILARITY))

	if got := len(req.GetSource().GetCollectionIds()); got != 2 {
		t.Fatalf("collection len mismatch: %d", got)
	}
	if req.GetLimit() != 5 {
		t.Fatalf("limit mismatch: %d", req.GetLimit())
	}
	if req.GetRankingMetric() != xaiapiv1.RankingMetric_RANKING_METRIC_COSINE_SIMILARITY {
		t.Fatalf("metric mismatch: %v", req.GetRankingMetric())
	}
}

func TestToolResponsePayload(t *testing.T) {
	matches := []*xaiapiv1.SearchMatch{
		{
			FileId:       "file1",
			ChunkId:      "chunk1",
			ChunkContent: "hello world",
			Score:        0.9,
		},
	}
	payload, err := ToolResponsePayload("call-1", matches, 1)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if payload == "" {
		t.Fatalf("payload empty")
	}
	msg, err := ToolResponseMessage("call-1", matches, 1)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if msg.GetRole() != xaiapiv1.MessageRole_ROLE_TOOL {
		t.Fatalf("role mismatch")
	}
}
