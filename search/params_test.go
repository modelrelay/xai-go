package search

import (
	"testing"
	"time"
)

func TestParametersValidation(t *testing.T) {
	src, err := WebSource(WebAllow("example.com"))
	if err != nil {
		t.Fatalf("web source err: %v", err)
	}
	p, err := Parameters(
		WithMode(2),
		WithSources(src),
		WithDateRange(time.Now().Add(-time.Hour), time.Now()),
		WithMaxResults(5),
	)
	if err != nil {
		t.Fatalf("params err: %v", err)
	}
	if len(p.GetSources()) != 1 {
		t.Fatalf("expected 1 source")
	}
}

func TestWebSourceValidation(t *testing.T) {
	_, err := WebSource(WebAllow("a.com"), WebExclude("b.com"))
	if err == nil {
		t.Fatalf("expected error for mutually exclusive domains")
	}
}

func TestRssSourceRequiresLinks(t *testing.T) {
	if _, err := RssSource(); err == nil {
		t.Fatalf("expected error for missing links")
	}
}
