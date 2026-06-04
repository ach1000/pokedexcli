package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/ach1000/pokedexcli/internal/pokeapi"
)

// commandMockHTTPClient is a test double for pokeapi.HTTPClient.
type commandMockHTTPClient struct {
	body       string
	statusCode int
	err        error
}

func (m *commandMockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(strings.NewReader(m.body)),
	}, nil
}

// captureStdout redirects os.Stdout for the duration of fn and returns what was printed.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func makeTestConfig(client pokeapi.HTTPClient) *config {
	return &config{
		nextLocationURL: pokeapi.LocationAreaURL,
		httpClient:      client,
	}
}

// --- commandMapBack tests ---

func TestCommandMapBack_FirstPage(t *testing.T) {
	cfg := makeTestConfig(nil) // prevLocationURL is "" — no HTTP call expected

	out := captureStdout(t, func() {
		if err := commandMapBack(cfg, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "you're on the first page") {
		t.Errorf("expected first-page message, got: %q", out)
	}
}

// --- commandMap tests ---

const twoPageLocationAreaJSON = `{
  "count": 2,
  "next": "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
  "previous": "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
  "results": [
    {"name": "bulbasaur-land", "url": "https://pokeapi.co/api/v2/location-area/99/"},
    {"name": "charmander-cave", "url": "https://pokeapi.co/api/v2/location-area/100/"}
  ]
}`

func TestCommandMap_PrintsAreaNames(t *testing.T) {
	client := &commandMockHTTPClient{body: twoPageLocationAreaJSON, statusCode: http.StatusOK}
	cfg := makeTestConfig(client)

	out := captureStdout(t, func() {
		if err := commandMap(cfg, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "bulbasaur-land") {
		t.Errorf("expected 'bulbasaur-land' in output, got: %q", out)
	}
	if !strings.Contains(out, "charmander-cave") {
		t.Errorf("expected 'charmander-cave' in output, got: %q", out)
	}
}

func TestCommandMap_UpdatesConfigURLs(t *testing.T) {
	client := &commandMockHTTPClient{body: twoPageLocationAreaJSON, statusCode: http.StatusOK}
	cfg := makeTestConfig(client)

	captureStdout(t, func() {
		commandMap(cfg, nil)
	})

	wantNext := "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20"
	wantPrev := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"

	if cfg.nextLocationURL != wantNext {
		t.Errorf("nextLocationURL: want %q, got %q", wantNext, cfg.nextLocationURL)
	}
	if cfg.prevLocationURL != wantPrev {
		t.Errorf("prevLocationURL: want %q, got %q", wantPrev, cfg.prevLocationURL)
	}
}

func TestCommandMap_PropagatesHTTPError(t *testing.T) {
	client := &commandMockHTTPClient{err: errors.New("network down")}
	cfg := makeTestConfig(client)

	err := commandMap(cfg, nil)
	if err == nil {
		t.Fatal("expected an error when HTTP client fails, got nil")
	}
}

// --- commandExplore tests ---

const exploreJSON = `{
  "pokemon_encounters": [
    {"pokemon": {"name": "tentacool", "url": ""}},
    {"pokemon": {"name": "gyarados",  "url": ""}}
  ]
}`

func TestCommandExplore_PrintsPokemon(t *testing.T) {
	client := &commandMockHTTPClient{body: exploreJSON, statusCode: http.StatusOK}
	cfg := makeTestConfig(client)

	out := captureStdout(t, func() {
		if err := commandExplore(cfg, []string{"pastoria-city-area"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "tentacool") {
		t.Errorf("expected 'tentacool' in output, got: %q", out)
	}
	if !strings.Contains(out, "gyarados") {
		t.Errorf("expected 'gyarados' in output, got: %q", out)
	}
}

func TestCommandExplore_MissingArgReturnsError(t *testing.T) {
	cfg := makeTestConfig(nil)

	err := commandExplore(cfg, nil)
	if err == nil {
		t.Fatal("expected error when no location arg given, got nil")
	}
}

func TestCommandExplore_PropagatesHTTPError(t *testing.T) {
	client := &commandMockHTTPClient{err: errors.New("network down")}
	cfg := makeTestConfig(client)

	err := commandExplore(cfg, []string{"pastoria-city-area"})
	if err == nil {
		t.Fatal("expected error when HTTP client fails, got nil")
	}
}
