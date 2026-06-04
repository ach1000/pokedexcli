package main

import (
	"strings"
	"testing"
)

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
