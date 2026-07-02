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
	//
	// Because chunks are typed, reasoning deltas and answer deltas are separate
	// fields — no parsing needed to tell them apart. Reasoning models like
	// grok-4.3 stream their thinking first; render it dimmed, then the answer.
	acc := responses.NewAccumulator()
	inReasoning := false
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
			delta := out.GetDelta()
			if r := delta.GetReasoningContent(); r != "" {
				if !inReasoning {
					fmt.Print("\x1b[2m") // dim the reasoning trace
					inReasoning = true
				}
				fmt.Print(r)
			}
			if c := delta.GetContent(); c != "" {
				if inReasoning {
					fmt.Print("\x1b[0m\n\n") // reset before the answer
					inReasoning = false
				}
				fmt.Print(c)
			}
		}
	}
	if inReasoning {
		fmt.Print("\x1b[0m")
	}

	final := acc.Response()
	if len(final.GetOutputs()) == 0 {
		log.Fatal("stream completed without producing any output")
	}
	fmt.Printf("\n\n✓ done — %d completion tokens; accumulator rebuilt the full typed response (model %s)\n",
		final.GetUsage().GetCompletionTokens(), final.GetModel())
}
