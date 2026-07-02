//go:build integration

package integration

import (
	"os"
	"testing"
)

func requireKey(t *testing.T) {
	t.Helper()
	if os.Getenv("XAI_API_KEY") == "" {
		t.Skip("set XAI_API_KEY to run integration tests")
	}
}
