package main

import (
	"bytes"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ach1000/pokedexcli/internal/pokeapi"
	"github.com/ach1000/pokedexcli/internal/pokecache"
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
	cache := pokecache.NewCache(1 * time.Minute)
	return &config{
		nextLocationURL: pokeapi.LocationAreaURL,
		httpClient:      client,
		cache:           cache,
		pokedex:         map[string]pokeapi.Pokemon{},
		randIntn:        rand.Intn,
	}
}

// --- commandHelp tests ---

func TestCommandHelp_PrintsKnownCommands(t *testing.T) {
	cfg := makeTestConfig(nil)

	out := captureStdout(t, func() {
		if err := commandHelp(cfg, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	for _, name := range []string{"help", "exit", "map", "mapb", "explore", "catch", "inspect", "pokedex", "cache"} {
		if !strings.Contains(out, name) {
			t.Errorf("expected %q in help output, got: %q", name, out)
		}
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

func TestCommandMapBack_ReturnsPreviousPage(t *testing.T) {
	client := &commandMockHTTPClient{body: twoPageLocationAreaJSON, statusCode: http.StatusOK}
	cfg := makeTestConfig(client)
	cfg.prevLocationURL = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"

	out := captureStdout(t, func() {
		if err := commandMapBack(cfg, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "bulbasaur-land") {
		t.Errorf("expected area names in output, got: %q", out)
	}
}

func TestCommandMapBack_PropagatesHTTPError(t *testing.T) {
	client := &commandMockHTTPClient{err: errors.New("network down")}
	cfg := makeTestConfig(client)
	cfg.prevLocationURL = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"

	err := commandMapBack(cfg, nil)
	if err == nil {
		t.Fatal("expected error when HTTP client fails, got nil")
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

// --- commandCatch tests ---

const catchPikachuJSON = `{
  "name": "pikachu",
  "base_experience": 112,
  "height": 4,
  "weight": 60,
  "stats": [],
  "types": []
}`

func TestCommandCatch_MissingArgReturnsError(t *testing.T) {
	cfg := makeTestConfig(nil)
	if err := commandCatch(cfg, nil); err == nil {
		t.Fatal("expected error when no pokemon arg given, got nil")
	}
}

func TestCommandCatch_PropagatesHTTPError(t *testing.T) {
	client := &commandMockHTTPClient{err: errors.New("network down")}
	cfg := makeTestConfig(client)
	if err := commandCatch(cfg, []string{"pikachu"}); err == nil {
		t.Fatal("expected error when HTTP client fails, got nil")
	}
}

func TestCommandCatch_CatchSuccess_AddsToPokedex(t *testing.T) {
	client := &commandMockHTTPClient{body: catchPikachuJSON, statusCode: http.StatusOK}
	cfg := makeTestConfig(client)
	// Always return 0 so 0 < 50 is always true → guaranteed catch.
	cfg.randIntn = func(_ int) int { return 0 }

	out := captureStdout(t, func() {
		if err := commandCatch(cfg, []string{"pikachu"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "was caught") {
		t.Errorf("expected 'was caught' in output, got: %q", out)
	}
	if !strings.Contains(out, "inspect command") {
		t.Errorf("expected inspect hint in output, got: %q", out)
	}
	if _, ok := cfg.pokedex["pikachu"]; !ok {
		t.Error("expected pikachu to be in pokedex after catch, but it wasn't")
	}
}

func TestCommandCatch_CatchFail_NotAddedToPokedex(t *testing.T) {
	client := &commandMockHTTPClient{body: catchPikachuJSON, statusCode: http.StatusOK}
	cfg := makeTestConfig(client)
	// Always return base_experience - 1 so result >= 50 → guaranteed escape.
	cfg.randIntn = func(n int) int { return n - 1 }

	out := captureStdout(t, func() {
		if err := commandCatch(cfg, []string{"pikachu"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "escaped") {
		t.Errorf("expected 'escaped' in output, got: %q", out)
	}
	if _, ok := cfg.pokedex["pikachu"]; ok {
		t.Error("expected pikachu NOT to be in pokedex after escape, but it was")
	}
}

// --- commandPokedex tests ---

func TestCommandPokedex_Empty(t *testing.T) {
	cfg := makeTestConfig(nil)

	out := captureStdout(t, func() {
		if err := commandPokedex(cfg, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "empty") {
		t.Errorf("expected empty-pokedex message, got: %q", out)
	}
}

func TestCommandPokedex_ListsCaughtPokemon(t *testing.T) {
	cfg := makeTestConfig(nil)
	cfg.pokedex["pikachu"] = pokeapi.Pokemon{Name: "pikachu"}
	cfg.pokedex["caterpie"] = pokeapi.Pokemon{Name: "caterpie"}

	out := captureStdout(t, func() {
		if err := commandPokedex(cfg, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "pikachu") {
		t.Errorf("expected 'pikachu' in output, got: %q", out)
	}
	if !strings.Contains(out, "caterpie") {
		t.Errorf("expected 'caterpie' in output, got: %q", out)
	}
}

// --- commandCache tests ---

func TestCommandCache_PrintsStats(t *testing.T) {
	cfg := makeTestConfig(nil)
	cfg.cache.Add("k1", []byte("v1"))

	out := captureStdout(t, func() {
		if err := commandCache(cfg, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "Cache items: 1") {
		t.Errorf("expected item count in output, got: %q", out)
	}
	if !strings.Contains(out, "Average lifetime:") {
		t.Errorf("expected average lifetime in output, got: %q", out)
	}
}

func TestCommandCache_NoConfiguredCache(t *testing.T) {
	cfg := makeTestConfig(nil)
	cfg.cache = nil

	out := captureStdout(t, func() {
		if err := commandCache(cfg, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "Cache is not configured") {
		t.Errorf("expected missing-cache message, got: %q", out)
	}
}

// --- commandInspect tests ---

const inspectPikachuJSON = `{
  "name": "pikachu",
  "base_experience": 112,
  "height": 4,
  "weight": 60,
  "stats": [
    {"base_stat": 35, "stat": {"name": "hp",     "url": ""}},
    {"base_stat": 55, "stat": {"name": "attack",  "url": ""}}
  ],
  "types": [
    {"type": {"name": "electric", "url": ""}}
  ]
}`

func TestCommandInspect_MissingArgReturnsError(t *testing.T) {
	cfg := makeTestConfig(nil)
	if err := commandInspect(cfg, nil); err == nil {
		t.Fatal("expected error when no pokemon arg given, got nil")
	}
}

func TestCommandInspect_NotCaught(t *testing.T) {
	cfg := makeTestConfig(nil)

	out := captureStdout(t, func() {
		if err := commandInspect(cfg, []string{"pikachu"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "you have not caught that pokemon") {
		t.Errorf("expected not-caught message, got: %q", out)
	}
}

func TestCommandInspect_PrintsDetails(t *testing.T) {
	client := &commandMockHTTPClient{body: catchPikachuJSON, statusCode: http.StatusOK}
	cfg := makeTestConfig(client)
	cfg.randIntn = func(_ int) int { return 0 } // guaranteed catch

	// Catch first so the pokemon is in the pokedex.
	// Use the richer inspectPikachuJSON for the inspect call.
	client.body = inspectPikachuJSON
	captureStdout(t, func() { commandCatch(cfg, []string{"pikachu"}) })

	out := captureStdout(t, func() {
		if err := commandInspect(cfg, []string{"pikachu"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	for _, want := range []string{"Name: pikachu", "Height: 4", "Weight: 60", "hp", "attack", "electric"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %q", want, out)
		}
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
