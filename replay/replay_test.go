package replay

import (
	"context"
	"net"
	"path/filepath"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

const bufSize = 1024 * 1024

type chatServer struct {
	xaiapiv1.UnimplementedChatServer
}

func (s *chatServer) GetCompletion(ctx context.Context, req *xaiapiv1.GetCompletionsRequest) (*xaiapiv1.GetChatCompletionResponse, error) {
	return &xaiapiv1.GetChatCompletionResponse{
		Id:    "resp_unary",
		Model: req.GetModel(),
		Outputs: []*xaiapiv1.CompletionOutput{{
			Index:   0,
			Message: &xaiapiv1.CompletionMessage{Content: "hello"},
		}},
	}, nil
}

func (s *chatServer) GetCompletionChunk(req *xaiapiv1.GetCompletionsRequest, stream xaiapiv1.Chat_GetCompletionChunkServer) error {
	chunks := []*xaiapiv1.GetChatCompletionChunk{
		{
			Id:    "resp_stream",
			Model: req.GetModel(),
			Outputs: []*xaiapiv1.CompletionOutputChunk{{
				Index: 0,
				Delta: &xaiapiv1.Delta{Content: "hi"},
			}},
		},
		{
			Id:    "resp_stream",
			Model: req.GetModel(),
			Outputs: []*xaiapiv1.CompletionOutputChunk{{
				Index: 0,
				Delta: &xaiapiv1.Delta{Content: " there"},
			}},
		},
	}
	for _, chunk := range chunks {
		if err := stream.Send(chunk); err != nil {
			return err
		}
	}
	return nil
}

func TestReplayHarness_RecordAndReplay(t *testing.T) {
	ctx := context.Background()
	listener := bufconn.Listen(bufSize)

	grpcServer := grpc.NewServer()
	xaiapiv1.RegisterChatServer(grpcServer, &chatServer{})
	go func() {
		_ = grpcServer.Serve(listener)
	}()
	defer grpcServer.Stop()

	dialer := func(ctx context.Context, _ string) (net.Conn, error) {
		return listener.Dial()
	}

	path := filepath.Join(t.TempDir(), "replay.ndjson")
	recorder, err := Open(path, ModeRecord)
	if err != nil {
		t.Fatalf("open recorder: %v", err)
	}
	defer func() {
		if err := recorder.Close(); err != nil {
			t.Fatalf("recorder close: %v", err)
		}
	}()

	recordConn, err := grpc.DialContext(
		ctx,
		"bufnet",
		append(
			recorder.DialOptions(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(dialer),
			grpc.WithBlock(),
		)...,
	)
	if err != nil {
		t.Fatalf("dial record: %v", err)
	}
	defer recordConn.Close()

	chatClient := xaiapiv1.NewChatClient(recordConn)
	request := &xaiapiv1.GetCompletionsRequest{Model: "grok-test"}

	unaryResp, err := chatClient.GetCompletion(ctx, request)
	if err != nil {
		t.Fatalf("unary call: %v", err)
	}
	if unaryResp.GetId() != "resp_unary" {
		t.Fatalf("unexpected unary id %q", unaryResp.GetId())
	}

	stream, err := chatClient.GetCompletionChunk(ctx, request)
	if err != nil {
		t.Fatalf("stream call: %v", err)
	}
	firstChunk, err := stream.Recv()
	if err != nil {
		t.Fatalf("stream recv: %v", err)
	}
	if firstChunk.GetId() != "resp_stream" {
		t.Fatalf("unexpected stream id %q", firstChunk.GetId())
	}
	for {
		_, recvErr := stream.Recv()
		if recvErr != nil {
			break
		}
	}

	replayer, err := Open(path, ModeReplay)
	if err != nil {
		t.Fatalf("open replay: %v", err)
	}
	defer func() {
		if err := replayer.Close(); err != nil {
			t.Fatalf("replayer close: %v", err)
		}
	}()

	replayConn, err := grpc.DialContext(
		ctx,
		"replay",
		append(
			replayer.DialOptions(),
			grpc.WithContextDialer(dialer),
			grpc.WithBlock(),
		)...,
	)
	if err != nil {
		t.Fatalf("dial replay: %v", err)
	}
	defer replayConn.Close()

	replayClient := xaiapiv1.NewChatClient(replayConn)
	replayUnary, err := replayClient.GetCompletion(ctx, request)
	if err != nil {
		t.Fatalf("replay unary: %v", err)
	}
	if replayUnary.GetId() != "resp_unary" {
		t.Fatalf("unexpected replay unary id %q", replayUnary.GetId())
	}

	replayStream, err := replayClient.GetCompletionChunk(ctx, request)
	if err != nil {
		t.Fatalf("replay stream: %v", err)
	}
	replayChunk, err := replayStream.Recv()
	if err != nil {
		t.Fatalf("replay recv: %v", err)
	}
	if replayChunk.GetId() != "resp_stream" {
		t.Fatalf("unexpected replay stream id %q", replayChunk.GetId())
	}
	for {
		_, recvErr := replayStream.Recv()
		if recvErr != nil {
			break
		}
	}
}
