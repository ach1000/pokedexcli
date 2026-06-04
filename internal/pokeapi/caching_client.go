package pokeapi

import (
	"bytes"
	"io"
	"net/http"

	"github.com/ach1000/pokedexcli/internal/pokecache"
)

// CachingClient wraps any HTTPClient and serves repeated requests from an
// in-memory cache, keyed by the request URL string.
type CachingClient struct {
	inner HTTPClient
	cache *pokecache.Cache
}

// NewCachingClient returns an HTTPClient that caches successful responses.
func NewCachingClient(inner HTTPClient, cache *pokecache.Cache) HTTPClient {
	return &CachingClient{inner: inner, cache: cache}
}

// Do checks the cache first. On a miss it delegates to the inner client,
// caches the response body, and returns a reconstructed response.
func (c *CachingClient) Do(req *http.Request) (*http.Response, error) {
	key := req.URL.String()

	if body, ok := c.cache.Get(key); ok {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}, nil
	}

	resp, err := c.inner.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// Only cache successful responses so errors are never persisted.
	if resp.StatusCode < http.StatusMultipleChoices {
		c.cache.Add(key, body)
	}

	resp.Body = io.NopCloser(bytes.NewReader(body))
	return resp, nil
}
