package responses

import xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"

// GeneratedTokens returns the generated text and reasoning token count.
// SamplingUsage.CompletionTokens counts text only, while ReasoningTokens is a
// disjoint counter. SamplingUsage.TotalTokens already includes both counters
// plus prompt tokens.
func GeneratedTokens(usage *xaiapiv1.SamplingUsage) int32 {
	return usage.GetCompletionTokens() + usage.GetReasoningTokens()
}
