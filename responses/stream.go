package responses

import (
	"context"
	"errors"
	"io"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

type chunkStream interface {
	Recv() (*xaiapiv1.GetChatCompletionChunk, error)
	CloseSend() error
}

// Stream wraps the gRPC stream that emits GetChatCompletionChunk messages.
type Stream struct {
	raw chunkStream
}

// Recv blocks until the next chunk is available or the stream ends.
func (s *Stream) Recv() (*xaiapiv1.GetChatCompletionChunk, error) {
	if s == nil || s.raw == nil {
		return nil, io.EOF
	}
	return s.raw.Recv()
}

// CloseSend closes the client side of the stream.
func (s *Stream) CloseSend() error {
	if s == nil || s.raw == nil {
		return nil
	}
	return s.raw.CloseSend()
}

// Accumulate drains the stream, returning a fully-hydrated completion response.
// If acc is nil a new accumulator is created automatically.
func (s *Stream) Accumulate(acc *Accumulator) (*xaiapiv1.GetChatCompletionResponse, error) {
	if acc == nil {
		acc = NewAccumulator()
	}
	for {
		chunk, err := s.Recv()
		if err == io.EOF {
			return acc.Response(), nil
		}
		if err != nil {
			return nil, err
		}
		acc.AddChunk(chunk)
	}
}

// Iterator returns a high-level iterator for streaming chunks.
func (s *Stream) Iterator(ctx context.Context) *Iterator {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Iterator{stream: s, ctx: ctx}
}

// ForEachChunk drains the stream and invokes fn for each chunk.
func (s *Stream) ForEachChunk(ctx context.Context, fn func(*xaiapiv1.GetChatCompletionChunk) error) error {
	it := s.Iterator(ctx)
	for {
		chunk, ok, err := it.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		if err := fn(chunk); err != nil {
			return err
		}
	}
}

// Iterator provides chunk-by-chunk iteration semantics.
type Iterator struct {
	stream *Stream
	ctx    context.Context
	err    error
}

// Next returns the next chunk. ok is false when iteration is complete.
func (it *Iterator) Next() (*xaiapiv1.GetChatCompletionChunk, bool, error) {
	if it.stream == nil {
		return nil, false, io.EOF
	}
	if it.err != nil {
		return nil, false, it.err
	}
	if err := it.ctx.Err(); err != nil {
		it.err = err
		_ = it.stream.CloseSend()
		return nil, false, err
	}
	chunk, err := it.stream.Recv()
	if errors.Is(err, io.EOF) {
		it.err = err
		return nil, false, io.EOF
	}
	if err != nil {
		it.err = err
		return nil, false, err
	}
	return chunk, true, nil
}
