package pokeapi

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ach1000/pokedexcli/internal/pokecache"
)

type countingHTTPClient struct {
	calls int
}

func (c *countingHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	c.calls++
	body := fmt.Sprintf("response-%d", c.calls)
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
	}, nil
}

func mustRequest(t *testing.T, rawURL string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		t.Fatalf("http.NewRequest(%q): %v", rawURL, err)
	}
	return req
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("io.ReadAll: %v", err)
	}
	return string(b)
}

func TestCachingClient_CanonicalizesDefaultLocationAreaPagination(t *testing.T) {
	inner := &countingHTTPClient{}
	cache := pokecache.NewCache(1 * time.Minute)
	client := NewCachingClient(inner, cache)

	resp1, err := client.Do(mustRequest(t, "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"))
	if err != nil {
		t.Fatalf("first Do returned error: %v", err)
	}
	firstBody := readBody(t, resp1)

	resp2, err := client.Do(mustRequest(t, "https://pokeapi.co/api/v2/location-area/"))
	if err != nil {
		t.Fatalf("second Do returned error: %v", err)
	}
	secondBody := readBody(t, resp2)

	if inner.calls != 1 {
		t.Fatalf("expected 1 inner HTTP call due to canonicalized cache key, got %d", inner.calls)
	}
	if firstBody != secondBody {
		t.Fatalf("expected cached body reuse, got first=%q second=%q", firstBody, secondBody)
	}
}

func TestCachingClient_DoesNotCanonicalizeNonPokeAPIHost(t *testing.T) {
	inner := &countingHTTPClient{}
	cache := pokecache.NewCache(1 * time.Minute)
	client := NewCachingClient(inner, cache)

	_, err := client.Do(mustRequest(t, "https://example.com/api/v2/location-area/?offset=0&limit=20"))
	if err != nil {
		t.Fatalf("first Do returned error: %v", err)
	}

	_, err = client.Do(mustRequest(t, "https://example.com/api/v2/location-area/"))
	if err != nil {
		t.Fatalf("second Do returned error: %v", err)
	}

	if inner.calls != 2 {
		t.Fatalf("expected 2 inner HTTP calls for non-PokeAPI host, got %d", inner.calls)
	}
}
