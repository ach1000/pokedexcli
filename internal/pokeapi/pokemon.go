package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const pokemonBaseURL = "https://pokeapi.co/api/v2/pokemon/"

// Pokemon holds the fields we care about from the Pokemon endpoint.
type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int          `json:"base_stat"`
		Stat     NamedResource `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type NamedResource `json:"type"`
	} `json:"types"`
}

// GetPokemon fetches a Pokemon by name or ID. Responses are served from the
// cache when available because client is a CachingClient in production.
func GetPokemon(name string, client HTTPClient) (Pokemon, error) {
	req, err := http.NewRequest(http.MethodGet, pokemonBaseURL+name, nil)
	if err != nil {
		return Pokemon{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusMultipleChoices {
		return Pokemon{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}

	var pokemon Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return Pokemon{}, err
	}

	return pokemon, nil
}
