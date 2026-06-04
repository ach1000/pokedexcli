package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/ach1000/pokedexcli/internal/pokeapi"
	"github.com/ach1000/pokedexcli/internal/pokecache"
)

type config struct {
	nextLocationURL string
	prevLocationURL string
	httpClient      pokeapi.HTTPClient
	cache           *pokecache.Cache
	pokedex         map[string]pokeapi.Pokemon
	randIntn        func(int) int
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all caught Pokemon",
			callback:    commandPokedex,
		},
		"cache": {
			name:        "cache",
			description: "Show cache stats",
			callback:    commandCache,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cache := pokecache.NewCache(5 * time.Minute)
	replConfig := &config{
		nextLocationURL: pokeapi.LocationAreaURL,
		httpClient:      pokeapi.NewCachingClient(&http.Client{}, cache),
		cache:           cache,
		pokedex:         map[string]pokeapi.Pokemon{},
		randIntn:        rand.Intn,
	}

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}

		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		command, ok := commands[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := command.callback(replConfig, words[1:]); err != nil {
			fmt.Println(err)
		}
	}
}
