package responses

import (
	"context"
	"errors"
	"time"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

var (
	// ErrDeferredExpired indicates the deferred result expired before completion.
	ErrDeferredExpired = errors.New("deferred completion expired")
	// ErrDeferredNoResponse indicates DONE status without a response payload.
	ErrDeferredNoResponse = errors.New("deferred completion returned no response")
)

type deferredFetcher func(context.Context) (*xaiapiv1.GetDeferredCompletionResponse, error)

func pollDeferred(ctx context.Context, interval time.Duration, fetch deferredFetcher) (*xaiapiv1.GetDeferredCompletionResponse, error) {
	current := interval
	maxInterval := 5 * time.Second
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(current):
			resp, err := fetch(ctx)
			if err != nil {
				return nil, err
			}
			switch resp.GetStatus() {
			case xaiapiv1.DeferredStatus_DONE:
				return resp, nil
			case xaiapiv1.DeferredStatus_EXPIRED:
				return nil, ErrDeferredExpired
			case xaiapiv1.DeferredStatus_PENDING:
				if current < maxInterval {
					current *= 2
					if current > maxInterval {
						current = maxInterval
					}
				}
				continue
			default:
				return nil, errors.New("unknown deferred status")
			}
		}
	}
}
