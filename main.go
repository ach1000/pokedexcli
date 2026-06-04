package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sort"
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

func commandExit(_ *config, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *config, _ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	preferredOrder := []string{"help", "exit", "map", "mapb", "explore", "catch", "inspect", "pokedex", "cache"}
	printed := map[string]struct{}{}

	for _, name := range preferredOrder {
		command, ok := commands[name]
		if !ok {
			continue
		}

		printed[name] = struct{}{}
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	commandNames := make([]string, 0, len(commands))
	for name := range commands {
		if _, alreadyPrinted := printed[name]; alreadyPrinted {
			continue
		}

		commandNames = append(commandNames, name)
	}

	sort.Strings(commandNames)

	for _, name := range commandNames {
		command := commands[name]
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	return nil
}

func commandMap(cfg *config, _ []string) error {
	locationAreas, err := pokeapi.GetLocationAreas(cfg.nextLocationURL, cfg.httpClient)
	if err != nil {
		return err
	}

	cfg.nextLocationURL = stringValue(locationAreas.Next)
	cfg.prevLocationURL = stringValue(locationAreas.Previous)

	for _, area := range locationAreas.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandMapBack(cfg *config, _ []string) error {
	if cfg.prevLocationURL == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	locationAreas, err := pokeapi.GetLocationAreas(cfg.prevLocationURL, cfg.httpClient)
	if err != nil {
		return err
	}

	cfg.nextLocationURL = stringValue(locationAreas.Next)
	cfg.prevLocationURL = stringValue(locationAreas.Previous)

	for _, area := range locationAreas.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandPokedex(cfg *config, _ []string) error {
	if len(cfg.pokedex) == 0 {
		fmt.Println("Your Pokedex is empty.")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for name := range cfg.pokedex {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

func commandCache(cfg *config, _ []string) error {
	if cfg.cache == nil {
		fmt.Println("Cache is not configured.")
		return nil
	}

	stats := cfg.cache.Stats()
	fmt.Printf("Cache items: %d\n", stats.ItemCount)
	fmt.Printf("Average lifetime: %s\n", stats.AverageLifetime.Round(time.Millisecond))
	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: inspect <pokemon>")
	}

	pokemon, ok := cfg.pokedex[args[0]]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, s := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: catch <pokemon>")
	}

	name := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	pokemon, err := pokeapi.GetPokemon(name, cfg.httpClient)
	if err != nil {
		return err
	}

	// Higher base_experience → harder to catch.
	// rand.Intn(base_experience) < 50 gives a reasonable spread:
	// easy (~40 exp): ~100%, mid (~112 exp): ~45%, legendary (~600 exp): ~8%.
	if cfg.randIntn(pokemon.BaseExperience) < 50 {
		cfg.pokedex[pokemon.Name] = pokemon
		fmt.Printf("%s was caught!\n", name)
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%s escaped!\n", name)
	}

	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: explore <location-area>")
	}

	name := args[0]
	fmt.Printf("Exploring %s...\n", name)

	result, err := pokeapi.ExploreLocationArea(name, cfg.httpClient)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, enc := range result.PokemonEncounters {
		fmt.Printf(" - %s\n", enc.Pokemon.Name)
	}

	return nil
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
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
