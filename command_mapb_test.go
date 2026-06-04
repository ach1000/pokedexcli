package main

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestCommandMapBack_FirstPage(t *testing.T) {
	cfg := makeTestConfig(nil)

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
