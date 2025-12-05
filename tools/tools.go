package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// FunctionTool builds a function tool definition. The schema parameter may be a
// JSON string/bytes or any Go value that can be marshaled to JSON.
func FunctionTool(name, description string, schema any) (*xaiapiv1.Tool, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("function name is required")
	}
	fn := &xaiapiv1.Function{
		Name:        name,
		Description: description,
	}
	if schema != nil {
		payload, err := encodeSchema(schema)
		if err != nil {
			return nil, err
		}
		fn.Parameters = payload
	}
	return &xaiapiv1.Tool{Tool: &xaiapiv1.Tool_Function{Function: fn}}, nil
}

// WebSearchOption mutates a WebSearch tool definition.
type WebSearchOption func(*xaiapiv1.WebSearch)

// WebSearchTool builds a WebSearch tool with the provided options.
func WebSearchTool(opts ...WebSearchOption) (*xaiapiv1.Tool, error) {
	cfg := &xaiapiv1.WebSearch{}
	for _, opt := range opts {
		opt(cfg)
	}
	if len(cfg.AllowedDomains) > 0 && len(cfg.ExcludedDomains) > 0 {
		return nil, errors.New("allowed_domains and excluded_domains are mutually exclusive")
	}
	return &xaiapiv1.Tool{Tool: &xaiapiv1.Tool_WebSearch{WebSearch: cfg}}, nil
}

// WithAllowedDomains restricts web search to the provided domains.
func WithAllowedDomains(domains ...string) WebSearchOption {
	return func(cfg *xaiapiv1.WebSearch) {
		cfg.AllowedDomains = append(cfg.AllowedDomains, clean(domains)...)
	}
}

// WithExcludedDomains omits results from the supplied domains.
func WithExcludedDomains(domains ...string) WebSearchOption {
	return func(cfg *xaiapiv1.WebSearch) {
		cfg.ExcludedDomains = append(cfg.ExcludedDomains, clean(domains)...)
	}
}

// EnableImageUnderstanding toggles vision support in downstream tools.
func EnableImageUnderstanding(enabled bool) WebSearchOption {
	return func(cfg *xaiapiv1.WebSearch) {
		cfg.EnableImageUnderstanding = protoBool(enabled)
	}
}

// XSearchOption mutates an XSearch tool definition.
type XSearchOption func(*xaiapiv1.XSearch) error

// XSearchTool builds an XSearch tool configuration.
func XSearchTool(opts ...XSearchOption) (*xaiapiv1.Tool, error) {
	cfg := &xaiapiv1.XSearch{}
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}
	if len(cfg.AllowedXHandles) > 0 && len(cfg.ExcludedXHandles) > 0 {
		return nil, errors.New("allowed_x_handles and excluded_x_handles are mutually exclusive")
	}
	return &xaiapiv1.Tool{Tool: &xaiapiv1.Tool_XSearch{XSearch: cfg}}, nil
}

// WithXDateRange limits X search to the provided interval.
func WithXDateRange(from, to time.Time) XSearchOption {
	return func(cfg *xaiapiv1.XSearch) error {
		if !from.IsZero() && !to.IsZero() && to.Before(from) {
			return errors.New("to date must be after from date")
		}
		if !from.IsZero() {
			cfg.FromDate = timestamppb.New(from)
		}
		if !to.IsZero() {
			cfg.ToDate = timestamppb.New(to)
		}
		return nil
	}
}

// AllowXHandles whitelists specific X handles (without '@').
func AllowXHandles(handles ...string) XSearchOption {
	return func(cfg *xaiapiv1.XSearch) error {
		cfg.AllowedXHandles = append(cfg.AllowedXHandles, clean(handles)...)
		return nil
	}
}

// ExcludeXHandles blacklists specific X handles.
func ExcludeXHandles(handles ...string) XSearchOption {
	return func(cfg *xaiapiv1.XSearch) error {
		cfg.ExcludedXHandles = append(cfg.ExcludedXHandles, clean(handles)...)
		return nil
	}
}

// EnableXImageUnderstanding toggles image understanding for X search.
func EnableXImageUnderstanding(enabled bool) XSearchOption {
	return func(cfg *xaiapiv1.XSearch) error {
		cfg.EnableImageUnderstanding = protoBool(enabled)
		return nil
	}
}

// EnableXVideoUnderstanding toggles video understanding for X search.
func EnableXVideoUnderstanding(enabled bool) XSearchOption {
	return func(cfg *xaiapiv1.XSearch) error {
		cfg.EnableVideoUnderstanding = protoBool(enabled)
		return nil
	}
}

// CodeExecutionTool enables the built-in code interpreter.
func CodeExecutionTool() *xaiapiv1.Tool {
	return &xaiapiv1.Tool{Tool: &xaiapiv1.Tool_CodeExecution{CodeExecution: &xaiapiv1.CodeExecution{}}}
}

// CollectionsSearchTool configures a server-side collections search.
func CollectionsSearchTool(collections []string, limit *int32) *xaiapiv1.Tool {
	cfg := &xaiapiv1.CollectionsSearch{
		CollectionIds: clean(collections),
	}
	if limit != nil {
		cfg.Limit = limit
	}
	return &xaiapiv1.Tool{Tool: &xaiapiv1.Tool_CollectionsSearch{CollectionsSearch: cfg}}
}

// DocumentSearchTool configures document search.
func DocumentSearchTool(limit *int32) *xaiapiv1.Tool {
	cfg := &xaiapiv1.DocumentSearch{Limit: limit}
	return &xaiapiv1.Tool{Tool: &xaiapiv1.Tool_DocumentSearch{DocumentSearch: cfg}}
}

// MCPOption mutates an MCP tool configuration.
type MCPOption func(*xaiapiv1.MCP)

// MCPTool registers a Model Context Protocol server.
func MCPTool(serverURL string, opts ...MCPOption) (*xaiapiv1.Tool, error) {
	serverURL = strings.TrimSpace(serverURL)
	if serverURL == "" {
		return nil, errors.New("serverURL is required for MCP tool")
	}
	cfg := &xaiapiv1.MCP{ServerUrl: serverURL}
	for _, opt := range opts {
		opt(cfg)
	}
	return &xaiapiv1.Tool{Tool: &xaiapiv1.Tool_Mcp{Mcp: cfg}}, nil
}

// WithServerLabel sets a label prefix for tool calls.
func WithServerLabel(label string) MCPOption {
	return func(cfg *xaiapiv1.MCP) {
		cfg.ServerLabel = strings.TrimSpace(label)
	}
}

// WithServerDescription sets an optional MCP description.
func WithServerDescription(desc string) MCPOption {
	return func(cfg *xaiapiv1.MCP) {
		cfg.ServerDescription = desc
	}
}

// WithAllowedToolNames restricts MCP tool names the model may call.
func WithAllowedToolNames(names ...string) MCPOption {
	return func(cfg *xaiapiv1.MCP) {
		cfg.AllowedToolNames = append(cfg.AllowedToolNames, clean(names)...)
	}
}

// WithAuthorization sets the Authorization header for MCP requests.
func WithAuthorization(token string) MCPOption {
	return func(cfg *xaiapiv1.MCP) {
		token = strings.TrimSpace(token)
		if token != "" {
			cfg.Authorization = protoString(token)
		}
	}
}

// WithBearerToken sets the Authorization header using the Bearer scheme.
func WithBearerToken(token string) MCPOption {
	return func(cfg *xaiapiv1.MCP) {
		token = strings.TrimSpace(token)
		if token != "" {
			cfg.Authorization = protoString(fmt.Sprintf("Bearer %s", token))
		}
	}
}

// WithExtraHeader adds a custom header to MCP requests.
func WithExtraHeader(key, value string) MCPOption {
	return func(cfg *xaiapiv1.MCP) {
		if cfg.ExtraHeaders == nil {
			cfg.ExtraHeaders = map[string]string{}
		}
		key = strings.TrimSpace(key)
		if key != "" {
			cfg.ExtraHeaders[key] = value
		}
	}
}

func clean(values []string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func protoBool(v bool) *bool {
	return &v
}

func protoString(v string) *string {
	return &v
}

func encodeSchema(schema any) (string, error) {
	switch v := schema.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("marshal schema: %w", err)
		}
		return string(b), nil
	}
}
