package responses

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// ToolCallEvent represents updates emitted while the model is constructing tool calls.
type ToolCallEvent struct {
	OutputIndex    int32
	CallID         string
	ArgumentsDelta string
	Complete       bool
	Call           *xaiapiv1.ToolCall
}

// ToolCallTracker collects tool call deltas across streamed chunks.
type ToolCallTracker struct {
	states map[string]*toolCallState
}

// NewToolCallTracker creates a tracker ready to consume streaming chunks.
func NewToolCallTracker() *ToolCallTracker {
	return &ToolCallTracker{
		states: make(map[string]*toolCallState),
	}
}

// ConsumeChunk processes a streamed chunk and returns tool call events emitted by this chunk.
func (t *ToolCallTracker) ConsumeChunk(chunk *xaiapiv1.GetChatCompletionChunk) []ToolCallEvent {
	if chunk == nil {
		return nil
	}
	var events []ToolCallEvent
	for _, outChunk := range chunk.GetOutputs() {
		idx := outChunk.GetIndex()
		delta := outChunk.GetDelta()
		if delta == nil {
			continue
		}
		for i, call := range delta.GetToolCalls() {
			if call == nil {
				continue
			}
			key := toolCallKey(call, idx, int32(i))
			state := t.ensureState(key, idx, call.GetId())
			events = append(events, state.apply(call))
		}
	}
	return events
}

// Get returns the current aggregated tool call by ID (if known).
func (t *ToolCallTracker) Get(id string) (*xaiapiv1.ToolCall, bool) {
	for _, state := range t.states {
		if state.id == id {
			return proto.Clone(state.call).(*xaiapiv1.ToolCall), true
		}
	}
	return nil, false
}

type toolCallState struct {
	key   string
	id    string
	index int32
	call  *xaiapiv1.ToolCall
	args  string
}

func (t *ToolCallTracker) ensureState(key string, idx int32, id string) *toolCallState {
	if state, ok := t.states[key]; ok {
		if id != "" {
			state.id = id
		}
		return state
	}
	state := &toolCallState{
		key:   key,
		id:    id,
		index: idx,
		call:  &xaiapiv1.ToolCall{},
	}
	t.states[key] = state
	return state
}

func (s *toolCallState) apply(delta *xaiapiv1.ToolCall) ToolCallEvent {
	if delta.GetId() != "" {
		s.id = delta.GetId()
	}
	if delta.GetType() != xaiapiv1.ToolCallType_TOOL_CALL_TYPE_CLIENT_SIDE_TOOL || s.call.GetType() == 0 {
		s.call.Type = delta.GetType()
	}
	if delta.GetStatus() != xaiapiv1.ToolCallStatus_TOOL_CALL_STATUS_IN_PROGRESS || s.call.GetStatus() == 0 {
		s.call.Status = delta.GetStatus()
	}
	if delta.GetErrorMessage() != "" {
		s.call.ErrorMessage = proto.String(delta.GetErrorMessage())
	}

	var argumentsDelta string
	if fn := delta.GetFunction(); fn != nil {
		target := s.ensureFunction()
		if fn.GetName() != "" {
			target.Name = fn.GetName()
		}
		if fn.GetArguments() != "" {
			argumentsDelta = fn.GetArguments()
			s.args += fn.GetArguments()
			target.Arguments = s.args
		}
	}

	complete := delta.GetStatus() == xaiapiv1.ToolCallStatus_TOOL_CALL_STATUS_COMPLETED ||
		delta.GetStatus() == xaiapiv1.ToolCallStatus_TOOL_CALL_STATUS_FAILED ||
		delta.GetStatus() == xaiapiv1.ToolCallStatus_TOOL_CALL_STATUS_INCOMPLETE

	return ToolCallEvent{
		OutputIndex:    s.index,
		CallID:         s.id,
		ArgumentsDelta: argumentsDelta,
		Complete:       complete,
		Call:           proto.Clone(s.call).(*xaiapiv1.ToolCall),
	}
}

func (s *toolCallState) ensureFunction() *xaiapiv1.FunctionCall {
	fn, ok := s.call.GetTool().(*xaiapiv1.ToolCall_Function)
	if !ok || fn.Function == nil {
		fn = &xaiapiv1.ToolCall_Function{
			Function: &xaiapiv1.FunctionCall{},
		}
		s.call.Tool = fn
	}
	return fn.Function
}

func toolCallKey(call *xaiapiv1.ToolCall, outputIdx int32, relativeIdx int32) string {
	if call.GetId() != "" {
		return call.GetId()
	}
	return fmt.Sprintf("%d:%d", outputIdx, relativeIdx)
}
