package replay

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// Mode controls whether the harness records or replays traffic.
type Mode string

const (
	ModeRecord Mode = "record"
	ModeReplay Mode = "replay"
)

const (
	kindUnary  = "unary"
	kindStream = "stream"
)

type entry struct {
	Kind      string   `json:"kind"`
	Method    string   `json:"method"`
	Request   string   `json:"request,omitempty"`
	Responses []string `json:"responses,omitempty"`
	Error     string   `json:"error,omitempty"`
}

// Harness records or replays gRPC calls using client interceptors.
type Harness struct {
	mode    Mode
	path    string
	mu      sync.Mutex
	entries []entry
	index   int
	file    *os.File
	writer  *bufio.Writer
}

// Open initializes a replay harness for the provided path and mode.
func Open(path string, mode Mode) (*Harness, error) {
	if path == "" {
		return nil, errors.New("replay: path required")
	}
	h := &Harness{
		mode: mode,
		path: path,
	}
	switch mode {
	case ModeRecord:
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		h.file = file
		h.writer = bufio.NewWriter(file)
	case ModeReplay:
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		if err := h.readEntries(file); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("replay: unknown mode")
	}
	return h, nil
}

// Close flushes any pending recordings and closes the underlying file.
func (h *Harness) Close() error {
	if h == nil || h.writer == nil {
		return nil
	}
	if err := h.writer.Flush(); err != nil {
		_ = h.file.Close()
		return err
	}
	return h.file.Close()
}

// DialOptions returns gRPC dial options to enable record/replay interceptors.
func (h *Harness) DialOptions() []grpc.DialOption {
	if h == nil {
		return nil
	}
	opts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(h.UnaryInterceptor()),
		grpc.WithChainStreamInterceptor(h.StreamInterceptor()),
	}
	if h.mode == ModeReplay {
		opts = append(opts,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(replayDialer),
		)
	}
	return opts
}

// UnaryInterceptor returns a gRPC unary interceptor for record/replay.
func (h *Harness) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if h == nil {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		switch h.mode {
		case ModeReplay:
			ent, err := h.nextEntry(kindUnary, method)
			if err != nil {
				return err
			}
			if err := h.validateRequest(ent.Request, req); err != nil {
				return err
			}
			if ent.Error != "" {
				return errors.New(ent.Error)
			}
			if len(ent.Responses) == 0 {
				return errors.New("replay: missing response payload")
			}
			return decodeMessage(ent.Responses[0], reply)
		case ModeRecord:
			ent := entry{Kind: kindUnary, Method: method}
			ent.Request = encodeRequest(req)
			err := invoker(ctx, method, req, reply, cc, opts...)
			if err != nil {
				ent.Error = err.Error()
				h.appendEntry(ent)
				return err
			}
			resp, encodeErr := encodeResponse(reply)
			if encodeErr != nil {
				return encodeErr
			}
			ent.Responses = []string{resp}
			h.appendEntry(ent)
			return nil
		default:
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}
}

// StreamInterceptor returns a gRPC streaming interceptor for record/replay.
func (h *Harness) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if h == nil {
			return streamer(ctx, desc, cc, method, opts...)
		}
		switch h.mode {
		case ModeReplay:
			ent, err := h.nextEntry(kindStream, method)
			if err != nil {
				return nil, err
			}
			return newReplayStream(ctx, ent), nil
		case ModeRecord:
			stream, err := streamer(ctx, desc, cc, method, opts...)
			if err != nil {
				h.appendEntry(entry{Kind: kindStream, Method: method, Error: err.Error()})
				return nil, err
			}
			return &recordingStream{
				ClientStream: stream,
				harness:      h,
				entry:        entry{Kind: kindStream, Method: method},
			}, nil
		default:
			return streamer(ctx, desc, cc, method, opts...)
		}
	}
}

func (h *Harness) readEntries(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var ent entry
		if err := json.Unmarshal(line, &ent); err != nil {
			return err
		}
		h.entries = append(h.entries, ent)
	}
	return scanner.Err()
}

func (h *Harness) nextEntry(kind, method string) (entry, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.index >= len(h.entries) {
		return entry{}, errors.New("replay: no more entries")
	}
	ent := h.entries[h.index]
	h.index++
	if ent.Kind != kind || ent.Method != method {
		return entry{}, errors.New("replay: entry mismatch")
	}
	return ent, nil
}

func (h *Harness) appendEntry(ent entry) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.writer == nil {
		return
	}
	data, err := json.Marshal(ent)
	if err != nil {
		return
	}
	_, _ = h.writer.Write(append(data, '\n'))
	_ = h.writer.Flush()
}

func (h *Harness) validateRequest(expected string, req interface{}) error {
	if expected == "" || req == nil {
		return nil
	}
	actual := encodeRequest(req)
	if actual == "" || actual == expected {
		return nil
	}
	return errors.New("replay: request payload mismatch")
}

func encodeRequest(req interface{}) string {
	return encodeMessageFrom(req)
}

func encodeResponse(resp interface{}) (string, error) {
	return encodeMessageFromWithErr(resp)
}

func encodeMessageFrom(msg interface{}) string {
	val, _ := encodeMessageFromWithErr(msg)
	return val
}

func encodeMessageFromWithErr(msg interface{}) (string, error) {
	protoMsg, ok := msg.(proto.Message)
	if !ok || protoMsg == nil {
		return "", nil
	}
	data, err := proto.Marshal(protoMsg)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func decodeMessage(payload string, msg interface{}) error {
	if payload == "" {
		return nil
	}
	protoMsg, ok := msg.(proto.Message)
	if !ok || protoMsg == nil {
		return errors.New("replay: response is not proto message")
	}
	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return err
	}
	return proto.Unmarshal(data, protoMsg)
}

type recordingStream struct {
	grpc.ClientStream
	harness   *Harness
	entry     entry
	finalized bool
}

func (s *recordingStream) SendMsg(m interface{}) error {
	if s.entry.Request == "" {
		s.entry.Request = encodeRequest(m)
	}
	return s.ClientStream.SendMsg(m)
}

func (s *recordingStream) RecvMsg(m interface{}) error {
	err := s.ClientStream.RecvMsg(m)
	if err == nil {
		if payload, encErr := encodeResponse(m); encErr == nil && payload != "" {
			s.entry.Responses = append(s.entry.Responses, payload)
		}
		return nil
	}
	if errors.Is(err, io.EOF) {
		s.finalize("")
		return err
	}
	s.finalize(err.Error())
	return err
}

func (s *recordingStream) finalize(errMsg string) {
	if s.finalized {
		return
	}
	s.finalized = true
	s.entry.Error = errMsg
	s.harness.appendEntry(s.entry)
}

type replayStream struct {
	ctx     context.Context
	entry   entry
	index   int
	headers metadata.MD
	trailer metadata.MD
}

func newReplayStream(ctx context.Context, ent entry) *replayStream {
	if ctx == nil {
		ctx = context.Background()
	}
	return &replayStream{
		ctx:     ctx,
		entry:   ent,
		headers: metadata.MD{},
		trailer: metadata.MD{},
	}
}

func (s *replayStream) Header() (metadata.MD, error) {
	return s.headers, nil
}

func (s *replayStream) Trailer() metadata.MD {
	return s.trailer
}

func (s *replayStream) CloseSend() error {
	return nil
}

func (s *replayStream) Context() context.Context {
	return s.ctx
}

func (s *replayStream) SendMsg(m interface{}) error {
	if s.entry.Request != "" {
		actual := encodeRequest(m)
		if actual != "" && actual != s.entry.Request {
			return errors.New("replay: request payload mismatch")
		}
	}
	return nil
}

func (s *replayStream) RecvMsg(m interface{}) error {
	if s.index >= len(s.entry.Responses) {
		if s.entry.Error != "" {
			return errors.New(s.entry.Error)
		}
		return io.EOF
	}
	payload := s.entry.Responses[s.index]
	s.index++
	return decodeMessage(payload, m)
}

func replayDialer(ctx context.Context, _ string) (net.Conn, error) {
	client, server := net.Pipe()
	go func() {
		defer server.Close()
		buf := make([]byte, 1024)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if _, err := server.Read(buf); err != nil {
				return
			}
		}
	}()
	return client, nil
}
