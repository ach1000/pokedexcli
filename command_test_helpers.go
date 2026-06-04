package main

import (
	"bytes"
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

const twoPageLocationAreaJSON = `{
  "count": 2,
  "next": "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
  "previous": "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
  "results": [
    {"name": "bulbasaur-land", "url": "https://pokeapi.co/api/v2/location-area/99/"},
    {"name": "charmander-cave", "url": "https://pokeapi.co/api/v2/location-area/100/"}
  ]
}`

const catchPikachuJSON = `{
  "name": "pikachu",
  "base_experience": 112,
  "height": 4,
  "weight": 60,
  "stats": [],
  "types": []
}`

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

const exploreJSON = `{
  "pokemon_encounters": [
    {"pokemon": {"name": "tentacool", "url": ""}},
    {"pokemon": {"name": "gyarados",  "url": ""}}
  ]
}`
