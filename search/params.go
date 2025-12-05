package search

import (
	"errors"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	xaiapiv1 "github.com/modelrelay/xai-go/gen/xai/api/v1"
)

// Parameters builds a SearchParameters struct with validation.
func Parameters(opts ...Option) (*xaiapiv1.SearchParameters, error) {
	sp := &xaiapiv1.SearchParameters{}
	for _, opt := range opts {
		if err := opt(sp); err != nil {
			return nil, err
		}
	}
	if sp.GetMaxSearchResults() < 0 {
		return nil, errors.New("max_search_results must be positive")
	}
	if sp.GetMode() == xaiapiv1.SearchMode_INVALID_SEARCH_MODE {
		sp.Mode = xaiapiv1.SearchMode_AUTO_SEARCH_MODE
	}
	return sp, nil
}

// Option mutates SearchParameters.
type Option func(*xaiapiv1.SearchParameters) error

// WithMode sets the search mode.
func WithMode(mode xaiapiv1.SearchMode) Option {
	return func(sp *xaiapiv1.SearchParameters) error {
		sp.Mode = mode
		return nil
	}
}

// WithSources appends search sources.
func WithSources(sources ...*xaiapiv1.Source) Option {
	return func(sp *xaiapiv1.SearchParameters) error {
		for _, src := range sources {
			if src != nil {
				sp.Sources = append(sp.Sources, src)
			}
		}
		return nil
	}
}

// WithDateRange configures the date window.
func WithDateRange(from, to time.Time) Option {
	return func(sp *xaiapiv1.SearchParameters) error {
		if !from.IsZero() && !to.IsZero() && to.Before(from) {
			return errors.New("to date must be after from date")
		}
		if !from.IsZero() {
			sp.FromDate = timestamppb.New(from)
		}
		if !to.IsZero() {
			sp.ToDate = timestamppb.New(to)
		}
		return nil
	}
}

// WithReturnCitations toggles citations output.
func WithReturnCitations(enabled bool) Option {
	return func(sp *xaiapiv1.SearchParameters) error {
		sp.ReturnCitations = enabled
		return nil
	}
}

// WithMaxResults limits results (1-30).
func WithMaxResults(limit int32) Option {
	return func(sp *xaiapiv1.SearchParameters) error {
		if limit < 1 || limit > 30 {
			return errors.New("max_results must be within [1,30]")
		}
		sp.MaxSearchResults = &limit
		return nil
	}
}

// WebSource builds a Source containing a WebSource payload.
func WebSource(opts ...WebSourceOption) (*xaiapiv1.Source, error) {
	cfg := &xaiapiv1.WebSource{}
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}
	if len(cfg.AllowedWebsites) > 0 && len(cfg.ExcludedWebsites) > 0 {
		return nil, errors.New("allowed_websites and excluded_websites are mutually exclusive")
	}
	return &xaiapiv1.Source{Source: &xaiapiv1.Source_Web{Web: cfg}}, nil
}

// NewsSource builds a Source for news search.
func NewsSource(opts ...NewsSourceOption) (*xaiapiv1.Source, error) {
	cfg := &xaiapiv1.NewsSource{}
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}
	if len(cfg.ExcludedWebsites) > 0 && cfg.GetCountry() == "" {
		// nothing special but keep API consistent
	}
	return &xaiapiv1.Source{Source: &xaiapiv1.Source_News{News: cfg}}, nil
}

// XSource builds a Source for X search.
func XSource(opts ...XSourceOption) (*xaiapiv1.Source, error) {
	cfg := &xaiapiv1.XSource{}
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}
	if len(cfg.IncludedXHandles) > 0 && len(cfg.ExcludedXHandles) > 0 {
		return nil, errors.New("included_x_handles and excluded_x_handles are mutually exclusive")
	}
	return &xaiapiv1.Source{Source: &xaiapiv1.Source_X{X: cfg}}, nil
}

// RssSource builds a Source for RSS feeds.
func RssSource(links ...string) (*xaiapiv1.Source, error) {
	cleaned := clean(links)
	if len(cleaned) == 0 {
		return nil, errors.New("rss source requires at least one link")
	}
	cfg := &xaiapiv1.RssSource{Links: cleaned}
	return &xaiapiv1.Source{Source: &xaiapiv1.Source_Rss{Rss: cfg}}, nil
}

// WebSourceOption mutates a WebSource.
type WebSourceOption func(*xaiapiv1.WebSource) error

func WebAllow(domains ...string) WebSourceOption {
	return func(ws *xaiapiv1.WebSource) error {
		ws.AllowedWebsites = append(ws.AllowedWebsites, clean(domains)...)
		return nil
	}
}

func WebExclude(domains ...string) WebSourceOption {
	return func(ws *xaiapiv1.WebSource) error {
		ws.ExcludedWebsites = append(ws.ExcludedWebsites, clean(domains)...)
		return nil
	}
}

func WebCountry(code string) WebSourceOption {
	return func(ws *xaiapiv1.WebSource) error {
		code = strings.ToUpper(strings.TrimSpace(code))
		if code != "" {
			ws.Country = protoString(code)
		}
		return nil
	}
}

func WebSafeSearch(enabled bool) WebSourceOption {
	return func(ws *xaiapiv1.WebSource) error {
		ws.SafeSearch = enabled
		return nil
	}
}

// NewsSourceOption mutates a NewsSource.
type NewsSourceOption func(*xaiapiv1.NewsSource) error

func NewsExclude(domains ...string) NewsSourceOption {
	return func(ns *xaiapiv1.NewsSource) error {
		ns.ExcludedWebsites = append(ns.ExcludedWebsites, clean(domains)...)
		return nil
	}
}

func NewsCountry(code string) NewsSourceOption {
	return func(ns *xaiapiv1.NewsSource) error {
		code = strings.ToUpper(strings.TrimSpace(code))
		if code != "" {
			ns.Country = protoString(code)
		}
		return nil
	}
}

func NewsSafeSearch(enabled bool) NewsSourceOption {
	return func(ns *xaiapiv1.NewsSource) error {
		ns.SafeSearch = enabled
		return nil
	}
}

// XSourceOption mutates an XSource.
type XSourceOption func(*xaiapiv1.XSource) error

func XInclude(handles ...string) XSourceOption {
	return func(xs *xaiapiv1.XSource) error {
		xs.IncludedXHandles = append(xs.IncludedXHandles, clean(handles)...)
		return nil
	}
}

func XExclude(handles ...string) XSourceOption {
	return func(xs *xaiapiv1.XSource) error {
		xs.ExcludedXHandles = append(xs.ExcludedXHandles, clean(handles)...)
		return nil
	}
}

func XPostFavoriteThreshold(count int32) XSourceOption {
	return func(xs *xaiapiv1.XSource) error {
		if count > 0 {
			xs.PostFavoriteCount = protoInt32(count)
		}
		return nil
	}
}

func XPostViewThreshold(count int32) XSourceOption {
	return func(xs *xaiapiv1.XSource) error {
		if count > 0 {
			xs.PostViewCount = protoInt32(count)
		}
		return nil
	}
}

func clean(values []string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func protoString(v string) *string { return &v }
func protoInt32(v int32) *int32    { return &v }
