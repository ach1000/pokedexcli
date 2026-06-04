package main

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

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
