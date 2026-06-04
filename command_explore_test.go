package main

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

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
