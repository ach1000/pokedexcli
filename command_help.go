package main

import (
	"fmt"
	"sort"
)

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
