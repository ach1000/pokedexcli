package main

import (
	"strings"
	"testing"
)

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
