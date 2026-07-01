package transport

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"google.golang.org/grpc/credentials"
)

type metadataCredentials struct {
	apiKey    string
	userAgent string

	cacheOnce sync.Once
	cache     map[string]string
}

func newMetadataCredentials(apiKey, userAgent string) credentials.PerRPCCredentials {
	return &metadataCredentials{
		apiKey:    strings.TrimSpace(apiKey),
		userAgent: strings.TrimSpace(userAgent),
	}
}

func (m *metadataCredentials) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	var err error
	m.cacheOnce.Do(func() {
		if m.apiKey == "" {
			err = fmt.Errorf("missing API key")
			return
		}
		meta := map[string]string{
			"authorization": fmt.Sprintf("Bearer %s", m.apiKey),
		}
		if m.userAgent != "" {
			meta["user-agent"] = m.userAgent
		}
		m.cache = meta
	})
	return m.cache, err
}

func (*metadataCredentials) RequireTransportSecurity() bool {
	return true
}
