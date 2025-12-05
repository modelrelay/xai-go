package config

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"

	xai "github.com/modelrelay/xai-go"
)

// Config is a declarative container for client configuration.
type Config struct {
	APIKey      string            `json:"api_key"`
	Address     string            `json:"address"`
	UserAgent   string            `json:"user_agent"`
	DefaultUser string            `json:"default_user"`
	Logger      *slog.Logger      `json:"-"`
	DialOptions []grpc.DialOption `json:"-"`
}

// Options converts the data map into xai.Option values.
func (c Config) Options() []xai.Option {
	var opts []xai.Option
	if c.APIKey != "" {
		opts = append(opts, xai.WithAPIKey(c.APIKey))
	}
	if c.Address != "" {
		opts = append(opts, xai.WithAddress(c.Address))
	}
	if c.UserAgent != "" {
		opts = append(opts, xai.WithUserAgent(c.UserAgent))
	}
	if c.DefaultUser != "" {
		opts = append(opts, xai.WithDefaultUser(c.DefaultUser))
	}
	if c.Logger != nil {
		opts = append(opts, xai.WithLogger(c.Logger))
	}
	if len(c.DialOptions) > 0 {
		opts = append(opts, xai.WithDialOptions(c.DialOptions...))
	}
	return opts
}

// NewClient creates a new xAI client using the declarative config.
func (c Config) NewClient(ctx context.Context, extra ...xai.Option) (*xai.Client, error) {
	opts := append(c.Options(), extra...)
	return xai.NewClient(ctx, opts...)
}
