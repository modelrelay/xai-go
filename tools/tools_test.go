package tools

import (
	"testing"
	"time"
)

func TestFunctionTool(t *testing.T) {
	tool, err := FunctionTool("weather", "Gets weather", map[string]any{"type": "object"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	fn := tool.GetFunction()
	if fn.GetName() != "weather" {
		t.Fatalf("name mismatch: %s", fn.GetName())
	}
	if fn.GetParameters() == "" {
		t.Fatalf("expected parameters to be set")
	}
}

func TestWebSearchValidation(t *testing.T) {
	_, err := WebSearchTool(WithAllowedDomains("a.com"), WithExcludedDomains("b.com"))
	if err == nil {
		t.Fatalf("expected error for mutually exclusive fields")
	}
}

func TestXSearchDateRange(t *testing.T) {
	now := time.Now().UTC()
	tool, err := XSearchTool(WithXDateRange(now.Add(-time.Hour), now))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if tool.GetXSearch().GetFromDate() == nil {
		t.Fatalf("expected from date")
	}
}

func TestMCPTool(t *testing.T) {
	tool, err := MCPTool("https://mcp.example.com", WithServerLabel("example"), WithAuthorization("token"))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if tool.GetMcp().GetServerLabel() != "example" {
		t.Fatalf("label mismatch")
	}
}

func TestWithBearerToken(t *testing.T) {
	tool, err := MCPTool("https://mcp.example.com", WithBearerToken("abc123"))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got := tool.GetMcp().GetAuthorization(); got != "Bearer abc123" {
		t.Fatalf("bearer mismatch: %s", got)
	}
}
