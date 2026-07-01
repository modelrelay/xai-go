package transport

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	// DefaultAddress is the public Grok/xAI gRPC endpoint.
	DefaultAddress = "api.x.ai:443"

	defaultUserAgent = "xai-go/0.1.1"
)

// Config holds dial parameters shared across all services.
type Config struct {
	Address     string
	APIKey      string
	UserAgent   string
	Logger      *slog.Logger
	DialOptions []grpc.DialOption
}

// Dial creates a gRPC connection configured with TLS, per-RPC auth metadata,
// and light-weight logging interceptors.
func Dial(ctx context.Context, cfg Config) (*grpc.ClientConn, error) {
	address := cfg.Address
	if address == "" {
		address = DefaultAddress
	}

	userAgent := cfg.UserAgent
	if userAgent == "" {
		userAgent = defaultUserAgent
	}

	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}

	creds := credentials.NewTLS(&tls.Config{})
	perRPC := newMetadataCredentials(cfg.APIKey, userAgent)

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithChainUnaryInterceptor(unaryLoggingInterceptor(logger)),
		grpc.WithChainStreamInterceptor(streamLoggingInterceptor(logger)),
	}
	dialOptions = append(dialOptions, cfg.DialOptions...)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, dialOptions...)
	if err != nil {
		return nil, fmt.Errorf("dial gRPC endpoint %s: %w", address, err)
	}

	return conn, nil
}
