package responses

import (
	"testing"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

func TestToolCallTrackerConsumesChunks(t *testing.T) {
	tracker := NewToolCallTracker()

	chunk := &xaiapiv1.GetChatCompletionChunk{
		Outputs: []*xaiapiv1.CompletionOutputChunk{
			{
				Index: 0,
				Delta: &xaiapiv1.Delta{
					ToolCalls: []*xaiapiv1.ToolCall{
						{
							Id:     "call_1",
							Status: xaiapiv1.ToolCallStatus_TOOL_CALL_STATUS_IN_PROGRESS,
							Tool: &xaiapiv1.ToolCall_Function{
								Function: &xaiapiv1.FunctionCall{
									Name:      "get_weather",
									Arguments: "{\"city\":",
								},
							},
						},
					},
				},
			},
		},
	}

	events := tracker.ConsumeChunk(chunk)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].ArgumentsDelta != "{\"city\":" {
		t.Fatalf("unexpected delta: %q", events[0].ArgumentsDelta)
	}
	if events[0].Complete {
		t.Fatalf("should not be complete yet")
	}

	chunk2 := &xaiapiv1.GetChatCompletionChunk{
		Outputs: []*xaiapiv1.CompletionOutputChunk{
			{
				Index: 0,
				Delta: &xaiapiv1.Delta{
					ToolCalls: []*xaiapiv1.ToolCall{
						{
							Id:     "call_1",
							Status: xaiapiv1.ToolCallStatus_TOOL_CALL_STATUS_COMPLETED,
							Tool: &xaiapiv1.ToolCall_Function{
								Function: &xaiapiv1.FunctionCall{
									Arguments: "\"SF\"}",
								},
							},
						},
					},
				},
			},
		},
	}

	events = tracker.ConsumeChunk(chunk2)
	if len(events) != 1 {
		t.Fatalf("expected 1 event on second chunk, got %d", len(events))
	}
	if !events[0].Complete {
		t.Fatalf("expected completion event")
	}
	if events[0].Call.GetFunction().GetArguments() != "{\"city\":\"SF\"}" {
		t.Fatalf("arguments mismatch: %s", events[0].Call.GetFunction().GetArguments())
	}
}
