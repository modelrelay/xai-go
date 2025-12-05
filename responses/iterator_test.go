package responses

import (
	"context"
	"errors"
	"io"
	"testing"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

type fakeChunkClient struct {
	chunks []*xaiapiv1.GetChatCompletionChunk
	err    error
	idx    int
	closed bool
}

func (f *fakeChunkClient) Recv() (*xaiapiv1.GetChatCompletionChunk, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.idx >= len(f.chunks) {
		return nil, io.EOF
	}
	ch := f.chunks[f.idx]
	f.idx++
	return ch, nil
}

func (f *fakeChunkClient) CloseSend() error {
	f.closed = true
	return nil
}

func TestIteratorNext(t *testing.T) {
	stream := &Stream{raw: &fakeChunkClient{
		chunks: []*xaiapiv1.GetChatCompletionChunk{
			{Id: "1"}, {Id: "2"},
		},
	}}
	it := stream.Iterator(context.Background())
	for i := 0; i < 2; i++ {
		chunk, ok, err := it.Next()
		if err != nil || !ok {
			t.Fatalf("unexpected err %v ok %v", err, ok)
		}
		if chunk.GetId() == "" {
			t.Fatalf("missing chunk id")
		}
	}
	if _, ok, err := it.Next(); !errors.Is(err, io.EOF) || ok {
		t.Fatalf("expected EOF, got err=%v ok=%v", err, ok)
	}
}

func TestIteratorContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	stream := &Stream{raw: &fakeChunkClient{}}
	it := stream.Iterator(ctx)
	cancel()
	if _, _, err := it.Next(); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancel, got %v", err)
	}
	fake := stream.raw.(*fakeChunkClient)
	if !fake.closed {
		t.Fatalf("expected stream closed after cancel")
	}
}

func TestForEachChunk(t *testing.T) {
	stream := &Stream{raw: &fakeChunkClient{
		chunks: []*xaiapiv1.GetChatCompletionChunk{{Id: "1"}},
	}}
	if err := stream.ForEachChunk(context.Background(), func(chunk *xaiapiv1.GetChatCompletionChunk) error {
		if chunk.GetId() != "1" {
			t.Fatalf("unexpected id %s", chunk.GetId())
		}
		return nil
	}); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}
