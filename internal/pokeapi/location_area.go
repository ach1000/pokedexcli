package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const LocationAreaURL = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"

type LocationAreaResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func GetLocationAreas(pageURL string) (LocationAreaResponse, error) {
	if pageURL == "" {
		pageURL = LocationAreaURL
	}

	resp, err := http.Get(pageURL)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusMultipleChoices {
		return LocationAreaResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	var locationAreas LocationAreaResponse
	if err := json.Unmarshal(body, &locationAreas); err != nil {
		return LocationAreaResponse{}, err
	}

	return locationAreas, nil
}
