package main

import (
	"net/http"
	"strings"
	"testing"
)

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
	cfg.randIntn = func(_ int) int { return 0 }

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
