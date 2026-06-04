# PROJECT.md

## Purpose
This file is a handoff reference for future coding agents working in this repository.

IMPORTANT: Update this file whenever behavior, structure, commands, assumptions, or key decisions change.

## Project Overview
- Language: Go
- Module: github.com/ach1000/pokedexcli
- Current app type: minimal CLI skeleton
- Current runtime behavior: prints `Hello, World!`

## Current File Map
- `main.go`: program entry point; prints a static greeting.
- `repl.go`: contains `cleanInput(text string) []string` utility.
- `repl_test.go`: table-driven tests for `cleanInput`.
- `Makefile`: convenience targets for build, run, test.

## cleanInput Contract
Function: `cleanInput(text string) []string`

Behavior:
1. Trims leading and trailing whitespace.
2. Converts input to lowercase.
3. Splits into words by any whitespace (spaces, tabs, newlines).
4. Returns an empty slice for empty or whitespace-only input.

Implementation details:
- Uses `strings.TrimSpace`.
- Uses `strings.ToLower`.
- Uses `strings.Fields` for tokenization.

## Test Coverage Notes
`repl_test.go` validates:
- Extra surrounding and interior spaces.
- Mixed uppercase/lowercase normalization.
- Tab/newline whitespace handling.
- Empty string input.

The tests compare both expected slice length and per-word values.

## Build and Execution Commands
Use either raw Go commands or Make targets.

Raw commands:
- `go build -o pokedexcli .`
- `go run .`
- `go test ./...`

Make targets:
- `make build`
- `make run`
- `make test`

## Assumptions and Constraints
- Package structure is currently single-package (`package main`).
- No external dependencies beyond Go standard library.
- No interactive REPL loop yet; `cleanInput` is prepared for future REPL parsing usage.
- Output binary name is `pokedexcli`.

## Suggested Next Evolutions
- Add actual REPL loop in `main.go` or a dedicated REPL module.
- Add command parsing and dispatch tests.
- Add error-path and edge-case tests once command handling exists.

## Maintenance Rule
When making further changes, update this file in the same PR/commit if any of the following are affected:
- Public behavior
- Commands or developer workflow
- Directory/file structure
- Assumptions or constraints
- Test strategy or guarantees
