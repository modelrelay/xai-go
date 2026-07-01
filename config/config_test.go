package config

import (
	"context"
	"testing"
)

func TestOptionsLength(t *testing.T) {
	cfg := Config{APIKey: "key", Address: "addr", UserAgent: "ua", DefaultUser: "user"}
	if got := len(cfg.Options()); got != 4 {
		t.Fatalf("unexpected option count: %d", got)
	}
}

func TestNewClientRequiresAPIKey(t *testing.T) {
	cfg := Config{}
	if _, err := cfg.NewClient(context.Background()); err == nil {
		t.Fatalf("expected error when api key missing")
	}
}
