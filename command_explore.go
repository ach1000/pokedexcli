package main

import (
	"fmt"

	"github.com/ach1000/pokedexcli/internal/pokeapi"
)

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
