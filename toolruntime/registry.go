package toolruntime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/messages"
	"github.com/modelrelay/xai-go/responses"
)

// Handler processes a tool call and returns an arbitrary payload that will be JSON encoded.
type Handler func(context.Context, *xaiapiv1.FunctionCall) (any, error)

// Registry maps function names to handlers to help drive tool orchestration loops.
type Registry struct {
	handlers map[string]Handler
}

// NewRegistry constructs an empty handler registry.
func NewRegistry() *Registry {
	return &Registry{handlers: map[string]Handler{}}
}

// Register adds or replaces a handler for the given function name.
func (r *Registry) Register(name string, handler Handler) {
	if r.handlers == nil {
		r.handlers = map[string]Handler{}
	}
	r.handlers[name] = handler
}

// Handle processes a completed tool call event and returns a ROLE_TOOL message containing the handler output.
func (r *Registry) Handle(ctx context.Context, event responses.ToolCallEvent) (*xaiapiv1.Message, error) {
	if !event.Complete {
		return nil, errors.New("tool call not completed")
	}
	call := event.Call.GetFunction()
	if call == nil {
		return nil, errors.New("tool call missing function payload")
	}
	handler, ok := r.handlers[call.GetName()]
	if !ok {
		return nil, fmt.Errorf("no handler registered for %s", call.GetName())
	}
	payload, err := handler(ctx, call)
	if err != nil {
		return nil, err
	}
	raw, err := encodeToolOutput(event.CallID, payload)
	if err != nil {
		return nil, err
	}
	callID := event.CallID
	msg := messages.SystemText("") // placeholder to reuse builders
	msg.Role = xaiapiv1.MessageRole_ROLE_TOOL
	msg.ToolCallId = &callID
	msg.Content = []*xaiapiv1.Content{
		{Content: &xaiapiv1.Content_Text{Text: raw}},
	}
	return msg, nil
}

func encodeToolOutput(callID string, payload any) (string, error) {
	body := map[string]any{
		"tool_call_id": callID,
		"output":       payload,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
