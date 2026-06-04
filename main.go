package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"github.com/ach1000/pokedexcli/internal/pokeapi"
)

type config struct {
	nextLocationURL string
	prevLocationURL string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
	}
}

func commandExit(_ *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	preferredOrder := []string{"help", "exit", "map", "mapb"}
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

func commandMap(cfg *config) error {
	locationAreas, err := pokeapi.GetLocationAreas(cfg.nextLocationURL)
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

func commandMapBack(cfg *config) error {
	if cfg.prevLocationURL == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	locationAreas, err := pokeapi.GetLocationAreas(cfg.prevLocationURL)
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

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	replConfig := &config{nextLocationURL: pokeapi.LocationAreaURL}

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

		if err := command.callback(replConfig); err != nil {
			fmt.Println(err)
		}
	}
}
