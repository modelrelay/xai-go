package responses

import (
	"sort"

	"google.golang.org/protobuf/proto"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// Accumulator rebuilds a GetChatCompletionResponse from streamed chunks.
type Accumulator struct {
	base    *xaiapiv1.GetChatCompletionResponse
	outputs map[int32]*xaiapiv1.CompletionOutput
}

// NewAccumulator creates an accumulator ready to ingest streaming chunks.
func NewAccumulator() *Accumulator {
	return &Accumulator{
		base:    &xaiapiv1.GetChatCompletionResponse{},
		outputs: map[int32]*xaiapiv1.CompletionOutput{},
	}
}

// AddChunk merges a streaming chunk into the aggregated response.
func (a *Accumulator) AddChunk(chunk *xaiapiv1.GetChatCompletionChunk) {
	if chunk == nil {
		return
	}
	if id := chunk.GetId(); id != "" {
		a.base.Id = id
	}
	if created := chunk.GetCreated(); created != nil {
		a.base.Created = created
	}
	if model := chunk.GetModel(); model != "" {
		a.base.Model = model
	}
	if fingerprint := chunk.GetSystemFingerprint(); fingerprint != "" {
		a.base.SystemFingerprint = fingerprint
	}
	if usage := chunk.GetUsage(); usage != nil {
		a.base.Usage = usage
	}
	if len(chunk.GetCitations()) > 0 {
		a.base.Citations = append([]string(nil), chunk.GetCitations()...)
	}
	if debug := chunk.GetDebugOutput(); debug != nil {
		a.base.DebugOutput = debug
	}

	for _, outChunk := range chunk.GetOutputs() {
		index := outChunk.GetIndex()
		output := a.ensureOutput(index)

		if delta := outChunk.GetDelta(); delta != nil {
			msg := output.GetMessage()
			if msg == nil {
				msg = &xaiapiv1.CompletionMessage{}
				output.Message = msg
			}
			msg.Content += delta.GetContent()
			msg.ReasoningContent += delta.GetReasoningContent()
			if delta.GetRole() != xaiapiv1.MessageRole_INVALID_ROLE {
				msg.Role = delta.GetRole()
			}
			if delta.GetEncryptedContent() != "" {
				msg.EncryptedContent += delta.GetEncryptedContent()
			}
			if len(delta.GetToolCalls()) > 0 {
				msg.ToolCalls = append(msg.ToolCalls, delta.GetToolCalls()...)
			}
		}

		if logprobs := outChunk.GetLogprobs(); logprobs != nil {
			dst := output.GetLogprobs()
			if dst == nil {
				dst = &xaiapiv1.LogProbs{}
				output.Logprobs = dst
			}
			dst.Content = append(dst.Content, logprobs.GetContent()...)
		}

		if reason := outChunk.GetFinishReason(); reason != xaiapiv1.FinishReason_REASON_INVALID {
			output.FinishReason = reason
		}
	}
}

// Response returns the aggregated completion response up to this point.
func (a *Accumulator) Response() *xaiapiv1.GetChatCompletionResponse {
	resp := proto.Clone(a.base).(*xaiapiv1.GetChatCompletionResponse)
	if len(a.outputs) == 0 {
		return resp
	}

	indexes := make([]int, 0, len(a.outputs))
	for idx := range a.outputs {
		indexes = append(indexes, int(idx))
	}
	sort.Ints(indexes)

	resp.Outputs = make([]*xaiapiv1.CompletionOutput, 0, len(indexes))
	for _, idx := range indexes {
		out := a.outputs[int32(idx)]
		resp.Outputs = append(resp.Outputs, proto.Clone(out).(*xaiapiv1.CompletionOutput))
	}
	return resp
}

func (a *Accumulator) ensureOutput(index int32) *xaiapiv1.CompletionOutput {
	if output, ok := a.outputs[index]; ok {
		return output
	}
	output := &xaiapiv1.CompletionOutput{
		Index:   index,
		Message: &xaiapiv1.CompletionMessage{},
	}
	a.outputs[index] = output
	return output
}
