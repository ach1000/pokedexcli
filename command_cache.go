package main

import (
	"fmt"
	"time"
)

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
