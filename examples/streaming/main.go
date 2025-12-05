package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	grok "github.com/modelrelay/xai-go"
	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/messages"
	"github.com/modelrelay/xai-go/responses"
)

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("GROK_API_KEY")
	if apiKey == "" {
		log.Fatal("set GROK_API_KEY before running the example")
	}
	client, err := grok.NewClient(ctx, grok.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("create client: %v", err)
	}
	defer client.Close()
	req := &xaiapiv1.GetCompletionsRequest{
		Model:    "grok-2-latest",
		Messages: []*xaiapiv1.Message{messages.UserText("Stream a haiku about databases.")},
	}
	stream, err := client.Responses.CreateStream(ctx, req)
	if err != nil {
		log.Fatalf("start stream: %v", err)
	}
	acc := responses.NewAccumulator()
	it := stream.Iterator(ctx)
	for {
		chunk, ok, err := it.Next()
		if errors.Is(err, io.EOF) || !ok {
			break
		}
		if err != nil {
			log.Fatalf("stream error: %v", err)
		}
		acc.AddChunk(chunk)
		for _, out := range chunk.GetOutputs() {
			fmt.Print(out.GetDelta().GetContent())
		}
	}
	full := acc.Response()
	fmt.Printf("\n\nFinal answer: %s\n", full.GetOutputs()[0].GetMessage().GetContent())
}
