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

// Reset clears the accumulator so it can be reused for another stream.
func (a *Accumulator) Reset() {
	a.base = &xaiapiv1.GetChatCompletionResponse{}
	for key := range a.outputs {
		delete(a.outputs, key)
	}
}

// AddChunk merges a streaming chunk into the aggregated response.
func (a *Accumulator) AddChunk(chunk *xaiapiv1.GetChatCompletionChunk) {
	mergeChunk(a.base, a.outputs, chunk)
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

// ReduceChunk returns a new response state by applying a chunk to the prior state.
// The input state is never mutated.
func ReduceChunk(state *xaiapiv1.GetChatCompletionResponse, chunk *xaiapiv1.GetChatCompletionChunk) *xaiapiv1.GetChatCompletionResponse {
	if state == nil && chunk == nil {
		return &xaiapiv1.GetChatCompletionResponse{}
	}
	var next *xaiapiv1.GetChatCompletionResponse
	if state == nil {
		next = &xaiapiv1.GetChatCompletionResponse{}
	} else {
		next = proto.Clone(state).(*xaiapiv1.GetChatCompletionResponse)
	}
	outputs := map[int32]*xaiapiv1.CompletionOutput{}
	for _, out := range next.Outputs {
		if out == nil {
			continue
		}
		outputs[out.Index] = out
	}
	mergeChunk(next, outputs, chunk)

	indexes := make([]int, 0, len(outputs))
	for idx := range outputs {
		indexes = append(indexes, int(idx))
	}
	sort.Ints(indexes)
	next.Outputs = make([]*xaiapiv1.CompletionOutput, 0, len(indexes))
	for _, idx := range indexes {
		next.Outputs = append(next.Outputs, outputs[int32(idx)])
	}
	return next
}

func mergeChunk(base *xaiapiv1.GetChatCompletionResponse, outputs map[int32]*xaiapiv1.CompletionOutput, chunk *xaiapiv1.GetChatCompletionChunk) {
	if chunk == nil || base == nil {
		return
	}
	if id := chunk.GetId(); id != "" {
		base.Id = id
	}
	if created := chunk.GetCreated(); created != nil {
		base.Created = created
	}
	if model := chunk.GetModel(); model != "" {
		base.Model = model
	}
	if fingerprint := chunk.GetSystemFingerprint(); fingerprint != "" {
		base.SystemFingerprint = fingerprint
	}
	if usage := chunk.GetUsage(); usage != nil {
		base.Usage = usage
	}
	if len(chunk.GetCitations()) > 0 {
		base.Citations = append([]string(nil), chunk.GetCitations()...)
	}
	if tier := chunk.GetServiceTier(); tier != xaiapiv1.ServiceTier_SERVICE_TIER_UNSPECIFIED {
		base.ServiceTier = tier
	}
	if files := chunk.GetOutputFiles(); len(files) > 0 {
		base.OutputFiles = append([]*xaiapiv1.OutputFile(nil), files...)
	}
	if debug := chunk.GetDebugOutput(); debug != nil {
		base.DebugOutput = debug
	}

	for _, outChunk := range chunk.GetOutputs() {
		index := outChunk.GetIndex()
		output := ensureOutput(outputs, index)

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
			if cites := delta.GetCitations(); len(cites) > 0 {
				msg.Citations = append(msg.Citations, cites...)
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

func (a *Accumulator) ensureOutput(index int32) *xaiapiv1.CompletionOutput {
	return ensureOutput(a.outputs, index)
}

func ensureOutput(outputs map[int32]*xaiapiv1.CompletionOutput, index int32) *xaiapiv1.CompletionOutput {
	if output, ok := outputs[index]; ok {
		return output
	}
	output := &xaiapiv1.CompletionOutput{
		Index:   index,
		Message: &xaiapiv1.CompletionMessage{},
	}
	outputs[index] = output
	return output
}
