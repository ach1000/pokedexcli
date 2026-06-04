package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const locationAreaBaseURL = "https://pokeapi.co/api/v2/location-area/"

// ExploreResponse holds the relevant fields from the location-area detail endpoint.
type ExploreResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

// PokemonEncounter is a single entry in the pokemon_encounters array.
type PokemonEncounter struct {
	Pokemon NamedResource `json:"pokemon"`
}

// NamedResource is a minimal named+URL pair used in several PokeAPI responses.
type NamedResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// ExploreLocationArea fetches the detail for a named location area and returns
// the list of Pokemon encounters. The response is served from the cache when
// available (because client is a CachingClient in production).
func ExploreLocationArea(name string, client HTTPClient) (ExploreResponse, error) {
	url := locationAreaBaseURL + name

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ExploreResponse{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return ExploreResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusMultipleChoices {
		return ExploreResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ExploreResponse{}, err
	}

	var result ExploreResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return ExploreResponse{}, err
	}

	return result, nil
}
