package main

import (
	"fmt"

	"github.com/ach1000/pokedexcli/internal/pokeapi"
)

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

	// Higher base_experience -> harder to catch.
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
