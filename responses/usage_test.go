package responses

import (
	"testing"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

func TestGeneratedTokens(t *testing.T) {
	tests := []struct {
		name  string
		usage *xaiapiv1.SamplingUsage
		want  int32
	}{
		{name: "nil usage"},
		{
			name: "text and reasoning",
			usage: &xaiapiv1.SamplingUsage{
				CompletionTokens: 1,
				ReasoningTokens:  1100,
				PromptTokens:     24,
				TotalTokens:      1125,
			},
			want: 1101,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GeneratedTokens(tt.usage); got != tt.want {
				t.Fatalf("GeneratedTokens() = %d, want %d", got, tt.want)
			}
		})
	}
}
