package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	grok "github.com/modelrelay/xai-go"
	"github.com/modelrelay/xai-go/documents"
	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/messages"
	"github.com/modelrelay/xai-go/responses"
	"github.com/modelrelay/xai-go/toolruntime"
	"github.com/modelrelay/xai-go/tools"
)

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("GROK_API_KEY")
	if apiKey == "" {
		log.Fatal("set GROK_API_KEY")
	}
	client, err := grok.NewClient(ctx, grok.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("client: %v", err)
	}
	defer client.Close()

	fnTool, err := tools.FunctionTool("lookup_docs", "Look up documents", map[string]any{
		"type":       "object",
		"properties": map[string]any{"query": map[string]any{"type": "string"}},
		"required":   []string{"query"},
	})
	if err != nil {
		log.Fatalf("function tool: %v", err)
	}
	req := &xaiapiv1.GetCompletionsRequest{
		Model: "grok-2-latest",
		Messages: []*xaiapiv1.Message{
			messages.SystemText("Answer with document search tool when needed."),
			messages.UserText("Summarize the onboarding guide for engineers."),
		},
		Tools: []*xaiapiv1.Tool{fnTool},
	}
	stream, err := client.Responses.CreateStream(ctx, req)
	if err != nil {
		log.Fatalf("stream: %v", err)
	}
	tracker := responses.NewToolCallTracker()
	registry := toolruntime.NewRegistry()
	registry.Register("lookup_docs", func(ctx context.Context, fn *xaiapiv1.FunctionCall) (any, error) {
		matchesResp, err := client.Documents.Search(ctx, documents.SearchRequest(fn.GetArguments(), documents.CollectionSource("engineering-guide")))
		if err != nil {
			return nil, err
		}
		return matchesResp.GetMatches(), nil
	})
	acc := responses.NewAccumulator()
	it := stream.Iterator(ctx)
	for {
		chunk, ok, err := it.Next()
		if errors.Is(err, io.EOF) || !ok {
			break
		}
		if err != nil {
			log.Fatalf("stream err: %v", err)
		}
		acc.AddChunk(chunk)
		for _, event := range tracker.ConsumeChunk(chunk) {
			if event.Complete {
				msg, err := registry.Handle(ctx, event)
				if err != nil {
					log.Fatalf("handler err: %v", err)
				}
				fmt.Printf("\nTool response: %s\n%s\n", event.CallID, msg.GetContent()[0].GetText())
			}
		}
	}
	full := acc.Response()
	fmt.Println("\nFinal answer:", full.GetOutputs()[0].GetMessage().GetContent())
}
