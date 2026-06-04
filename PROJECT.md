# PROJECT.md

## Purpose
This file is a handoff reference for future coding agents working in this repository.

IMPORTANT: Update this file whenever behavior, structure, commands, assumptions, or key decisions change.

## Project Overview
- Language: Go
- Module: github.com/ach1000/pokedexcli
- Current app type: interactive CLI with simple REPL loop
- Current runtime behavior: prompts with `Pokedex > `, dispatches commands via a registry, supports `help`, `exit`, `map`, `mapb`, `explore`, `catch`, and `inspect`

## Current File Map
- `main.go`: program entry point; runs an infinite REPL loop, defines command registry, and dispatches callbacks with shared REPL config.
- `internal/pokeapi/location_area.go`: PokeAPI client for paginated location-area requests; accepts an `HTTPClient` interface for testability.
- `internal/pokeapi/explore.go`: `ExploreLocationArea` — fetches Pokemon encounter list for a named location area.
- `internal/pokeapi/pokemon.go`: `GetPokemon` — fetches a single Pokemon by name; `Pokemon` struct with name, base_experience, height, weight, stats, types.
- `internal/pokeapi/caching_client.go`: `CachingClient` wraps any `HTTPClient` + `*pokecache.Cache`; serves cached responses on hit, caches 2xx responses on miss.
- `internal/pokecache/cache.go`: thread-safe in-memory cache with configurable TTL and background reap loop.
- `internal/pokecache/cache_test.go`: unit tests for Add/Get, miss, overwrite, reap eviction, and reap preservation.
- `commands_test.go`: unit tests for all command handlers using a mock HTTP client.
- `repl.go`: contains `cleanInput(text string) []string` utility.
- `repl_test.go`: table-driven tests for `cleanInput`.
- `Makefile`: convenience targets for build, run, test, clean.

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

`internal/pokecache/cache_test.go` validates:
- Add/Get round-trip, cache miss, overwrite semantics.
- Reap loop evicts entries older than the interval.
- Reap loop preserves entries that were added recently.

`internal/pokeapi/location_area_test.go` validates:
- Successful JSON parsing, default URL fallback, non-2xx error, HTTP error, invalid JSON.

`commands_test.go` validates:
- `commandHelp` lists all registered command names.
- `commandMapBack` first-page guard, happy-path pagination, HTTP error propagation.
- `commandMap` prints area names, updates config URLs, propagates HTTP errors.
- `commandCatch` missing arg, HTTP error, guaranteed-catch adds to pokedex, guaranteed-escape does not.
- `commandInspect` missing arg, not-yet-caught message, full detail output (name/height/weight/stats/types).
- `commandExplore` prints Pokemon names, missing arg, HTTP error.

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
- `make clean`

`make clean` removes generated output artifacts:
- `pokedexcli` (built binary)
- `repl.log` (captured CLI output log)

## Assumptions and Constraints
- Package structure: `main` package plus `internal/pokeapi` and `internal/pokecache`.
- No external dependencies beyond Go standard library.
- REPL uses the first normalized token as the command key in a command registry.
- Command callbacks share state via `*config` (next/previous location URLs, `HTTPClient`, `pokedex map[string]Pokemon`, `randIntn func(int) int`).
- The `HTTPClient` stored in `config` is always a `*pokeapi.CachingClient` wrapping `*http.Client` and a 5-minute `pokecache.Cache`.
- `randIntn` is injectable (defaults to `rand.Intn`) so catch-command tests can deterministically force catch or escape.
- Output binary name is `pokedexcli`.

## REPL Runtime Behavior
- Prompt: `Pokedex > ` (no newline before input).
- Input read via `bufio.Scanner` from `os.Stdin`.
- Input is normalized via `cleanInput`.
- Empty input is ignored.
- Supported commands:
	- `help`: prints a welcome message and dynamically lists registered commands.
	- `exit`: prints `Closing the Pokedex... Goodbye!` and exits the program.
	- `map`: fetches and prints the next 20 location-area names from PokeAPI.
	- `mapb`: fetches and prints the previous 20 location-area names.
	- `explore <area>`: lists all Pokemon encounters in a named location area.
	- `catch <pokemon>`: attempts to catch a Pokemon; chance based on base_experience; adds to pokedex on success.
	- `inspect <pokemon>`: prints name, height, weight, stats, and types for a caught Pokemon.
- `mapb` on the first page prints `you're on the first page`.
- `inspect` on an uncaught Pokemon prints `you have not caught that pokemon`.
- All commands take `words[1:]` as an args slice; no-arg commands ignore it.

## Suggested Next Evolutions
- Move command callbacks into dedicated files as command count grows.
- Add a `pokedex` command to list all caught Pokemon.

## Maintenance Rule
When making further changes, update this file in the same PR/commit if any of the following are affected:
- Public behavior
- Commands or developer workflow
- Directory/file structure
- Assumptions or constraints
- Test strategy or guarantees
