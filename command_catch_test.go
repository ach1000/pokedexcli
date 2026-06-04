package main

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

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
