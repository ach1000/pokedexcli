package main

import (
	"fmt"

	"github.com/ach1000/pokedexcli/internal/pokeapi"
)

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

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
