package main

import "strings"

func cleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return []string{}
	}

	lower := strings.ToLower(trimmed)
	return strings.Fields(lower)
}
