package toolruntime

import (
	"context"
	"testing"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/responses"
)

func TestRegistryHandle(t *testing.T) {
	reg := NewRegistry()
	reg.Register("echo", func(ctx context.Context, fn *xaiapiv1.FunctionCall) (any, error) {
		return map[string]string{"args": fn.GetArguments()}, nil
	})

	event := responses.ToolCallEvent{
		CallID:   "tool_1",
		Complete: true,
		Call: &xaiapiv1.ToolCall{
			Tool: &xaiapiv1.ToolCall_Function{Function: &xaiapiv1.FunctionCall{
				Name:      "echo",
				Arguments: "{\"text\":\"hi\"}",
			}},
		},
	}

	msg, err := reg.Handle(context.Background(), event)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if msg.GetRole() != xaiapiv1.MessageRole_ROLE_TOOL {
		t.Fatalf("role mismatch")
	}
}
