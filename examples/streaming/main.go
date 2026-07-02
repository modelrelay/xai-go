package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/modelrelay/xai-go"
	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
	"github.com/modelrelay/xai-go/messages"
	"github.com/modelrelay/xai-go/responses"
)

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		log.Fatal("set XAI_API_KEY before running the example")
	}
	client, err := xai.NewClient(ctx, xai.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("create client: %v", err)
	}
	defer client.Close()

	stream, err := client.Responses.CreateStream(ctx, &xaiapiv1.GetCompletionsRequest{
		Model:    "grok-4.3",
		Messages: []*xaiapiv1.Message{messages.UserText("Write an eight-line poem about streaming data.")},
	})
	if err != nil {
		log.Fatalf("start stream: %v", err)
	}

	// Drain the stream, printing tokens as they arrive and accumulating the full
	// response. Recv returns io.EOF at the end; any other error is a real gRPC
	// status and must not be swallowed.
	acc := responses.NewAccumulator()
	for {
		chunk, err := stream.Recv()
		if errors.Is(err, io.EOF) {
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

	outs := acc.Response().GetOutputs()
	if len(outs) == 0 {
		log.Fatal("stream completed without producing any output")
	}
	fmt.Printf("\n\nFinal answer: %s\n", outs[0].GetMessage().GetContent())
}
