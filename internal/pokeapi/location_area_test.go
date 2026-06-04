package pokeapi

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPClient is a test double that returns a preset response or error.
type mockHTTPClient struct {
	body       string
	statusCode int
	err        error
	requestURL string // populated on each Do call
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.requestURL = req.URL.String()
	if m.err != nil {
		return nil, m.err
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(strings.NewReader(m.body)),
	}, nil
}

const sampleLocationAreaJSON = `{
  "count": 3,
  "next": "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
  "previous": null,
  "results": [
    {"name": "canalave-city-area", "url": "https://pokeapi.co/api/v2/location-area/1/"},
    {"name": "eterna-city-area",   "url": "https://pokeapi.co/api/v2/location-area/2/"},
    {"name": "pastoria-city-area", "url": "https://pokeapi.co/api/v2/location-area/3/"}
  ]
}`

func TestGetLocationAreas_Success(t *testing.T) {
	client := &mockHTTPClient{body: sampleLocationAreaJSON, statusCode: http.StatusOK}

	result, err := GetLocationAreas("", client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Count != 3 {
		t.Errorf("expected count 3, got %d", result.Count)
	}
	if len(result.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result.Results))
	}
	if result.Results[0].Name != "canalave-city-area" {
		t.Errorf("unexpected first area name: %q", result.Results[0].Name)
	}
	if result.Previous != nil {
		t.Errorf("expected previous to be nil, got %v", result.Previous)
	}
	if result.Next == nil || *result.Next != "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20" {
		t.Errorf("unexpected next URL: %v", result.Next)
	}
}

func TestGetLocationAreas_FallsBackToDefaultURL(t *testing.T) {
	client := &mockHTTPClient{body: sampleLocationAreaJSON, statusCode: http.StatusOK}

	_, err := GetLocationAreas("", client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.requestURL != LocationAreaURL {
		t.Errorf("expected default URL %q, got %q", LocationAreaURL, client.requestURL)
	}
}

func TestGetLocationAreas_ErrorStatus(t *testing.T) {
	client := &mockHTTPClient{body: "not found", statusCode: http.StatusNotFound}

	_, err := GetLocationAreas("", client)
	if err == nil {
		t.Fatal("expected an error for non-2xx status, got nil")
	}
}

func TestGetLocationAreas_HTTPError(t *testing.T) {
	client := &mockHTTPClient{err: errors.New("connection refused")}

	_, err := GetLocationAreas("", client)
	if err == nil {
		t.Fatal("expected an error when HTTP client fails, got nil")
	}
}

func TestGetLocationAreas_InvalidJSON(t *testing.T) {
	client := &mockHTTPClient{body: `{invalid json`, statusCode: http.StatusOK}

	_, err := GetLocationAreas("", client)
	if err == nil {
		t.Fatal("expected an error for invalid JSON, got nil")
	}
}
