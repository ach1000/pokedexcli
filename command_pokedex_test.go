package main

import (
	"strings"
	"testing"

	"github.com/ach1000/pokedexcli/internal/pokeapi"
)

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
