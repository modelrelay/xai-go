package grok

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"google.golang.org/grpc"

	"github.com/modelrelay/xai-go/auth"
	"github.com/modelrelay/xai-go/batch"
	"github.com/modelrelay/xai-go/chat"
	"github.com/modelrelay/xai-go/documents"
	"github.com/modelrelay/xai-go/embeddings"
	"github.com/modelrelay/xai-go/files"
	"github.com/modelrelay/xai-go/images"
	"github.com/modelrelay/xai-go/internal/raw"
	"github.com/modelrelay/xai-go/internal/transport"
	"github.com/modelrelay/xai-go/models"
	"github.com/modelrelay/xai-go/responses"
	"github.com/modelrelay/xai-go/sample"
	"github.com/modelrelay/xai-go/tokenize"
	"github.com/modelrelay/xai-go/video"
)

const (
	// Version identifies the current SDK revision reported via User-Agent.
	Version = "0.1.0-dev"

	envAPIKey  = "GROK_API_KEY"
	envAddress = "GROK_GRPC_ADDRESS"
)

// Client holds service entry points for interacting with the Grok API.
type Client struct {
	conn *grpc.ClientConn

	Chat       chat.Service
	Responses  responses.Service
	Sample     sample.Service
	Embeddings embeddings.Service
	Images     images.Service
	Documents  documents.Service
	Tokenize   tokenize.Service
	Models     models.Service
	Auth       auth.Service
	Batch      batch.Service
	Files      files.Service
	Video      video.Service
}

// Option configures a client instance.
type Option func(*config)

type config struct {
	APIKey      string
	Address     string
	UserAgent   string
	DefaultUser string
	Logger      *slog.Logger
	DialOptions []grpc.DialOption
}

// NewClient constructs a new SDK client with the supplied options.
func NewClient(ctx context.Context, opts ...Option) (*Client, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("missing API key (set GROK_API_KEY or use WithAPIKey)")
	}

	conn, err := transport.Dial(ctx, transport.Config{
		Address:     cfg.Address,
		APIKey:      cfg.APIKey,
		UserAgent:   cfg.UserAgent,
		Logger:      cfg.Logger,
		DialOptions: cfg.DialOptions,
	})
	if err != nil {
		return nil, err
	}

	defaults := raw.Defaults{User: cfg.DefaultUser}
	rawClients := raw.NewClients(conn, defaults)

	chatService := chat.NewService(rawClients.Chat)
	client := &Client{
		conn:       conn,
		Chat:       chatService,
		Responses:  responses.NewService(chatService),
		Sample:     sample.NewService(rawClients.Sample),
		Embeddings: embeddings.NewService(rawClients.Embed),
		Images:     images.NewService(rawClients.Image),
		Documents:  documents.NewService(rawClients.Docs),
		Tokenize:   tokenize.NewService(rawClients.Token),
		Models:     models.NewService(rawClients.Models),
		Auth:       auth.NewService(rawClients.Auth),
		Batch:      batch.NewService(rawClients.Batch),
		Files:      files.NewService(rawClients.Files),
		Video:      video.NewService(rawClients.Video),
	}

	return client, nil
}

// Close cleans up the underlying gRPC connection.
func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// WithAPIKey sets the API key used for authentication.
func WithAPIKey(key string) Option {
	return func(cfg *config) {
		cfg.APIKey = strings.TrimSpace(key)
	}
}

// WithAddress overrides the default gRPC endpoint.
func WithAddress(address string) Option {
	return func(cfg *config) {
		if address != "" {
			cfg.Address = address
		}
	}
}

// WithUserAgent overrides the default user-agent header.
func WithUserAgent(userAgent string) Option {
	return func(cfg *config) {
		if userAgent != "" {
			cfg.UserAgent = userAgent
		}
	}
}

// WithDefaultUser sets a default request user identifier for Chat/Sample APIs.
func WithDefaultUser(user string) Option {
	return func(cfg *config) {
		cfg.DefaultUser = strings.TrimSpace(user)
	}
}

// WithLogger overrides the logger used by transport interceptors.
func WithLogger(logger *slog.Logger) Option {
	return func(cfg *config) {
		if logger != nil {
			cfg.Logger = logger
		}
	}
}

// WithDialOptions appends custom grpc.DialOption values.
func WithDialOptions(opts ...grpc.DialOption) Option {
	return func(cfg *config) {
		cfg.DialOptions = append(cfg.DialOptions, opts...)
	}
}

func defaultConfig() config {
	cfg := config{
		Address:   transport.DefaultAddress,
		UserAgent: fmt.Sprintf("grok-go-sdk/%s", Version),
		Logger:    slog.Default(),
	}

	if v := strings.TrimSpace(os.Getenv(envAPIKey)); v != "" {
		cfg.APIKey = v
	}
	if v := strings.TrimSpace(os.Getenv(envAddress)); v != "" {
		cfg.Address = v
	}

	return cfg
}
