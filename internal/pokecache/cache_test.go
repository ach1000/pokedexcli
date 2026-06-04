package pokecache

import (
	"testing"
	"time"
)

func TestCache_AddAndGet(t *testing.T) {
	c := NewCache(5 * time.Minute)

	c.Add("key1", []byte("value1"))
	c.Add("key2", []byte("value2"))

	got, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected key1 to be present")
	}
	if string(got) != "value1" {
		t.Errorf("key1: want %q, got %q", "value1", string(got))
	}

	got, ok = c.Get("key2")
	if !ok {
		t.Fatal("expected key2 to be present")
	}
	if string(got) != "value2" {
		t.Errorf("key2: want %q, got %q", "value2", string(got))
	}
}

func TestCache_MissReturnsFalse(t *testing.T) {
	c := NewCache(5 * time.Minute)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("expected miss for nonexistent key, got hit")
	}
}

func TestCache_OverwriteEntry(t *testing.T) {
	c := NewCache(5 * time.Minute)
	c.Add("key", []byte("first"))
	c.Add("key", []byte("second"))

	got, ok := c.Get("key")
	if !ok {
		t.Fatal("expected key to be present after overwrite")
	}
	if string(got) != "second" {
		t.Errorf("want %q after overwrite, got %q", "second", string(got))
	}
}

func TestCache_ReapRemovesExpiredEntries(t *testing.T) {
	interval := 50 * time.Millisecond
	c := NewCache(interval)

	c.Add("stale", []byte("data"))

	// Wait for the reap loop to fire at least once after the entry has aged past interval.
	time.Sleep(interval * 3)

	_, ok := c.Get("stale")
	if ok {
		t.Error("expected stale entry to be reaped, but it still exists")
	}
}

func TestCache_ReapPreservesNewEntries(t *testing.T) {
	interval := 100 * time.Millisecond
	c := NewCache(interval)

	// Sleep until just before the reap window, then add a fresh entry.
	time.Sleep(interval / 2)
	c.Add("fresh", []byte("keep me"))

	// Let one reap cycle run — old entries would be evicted, but "fresh" should survive.
	time.Sleep(interval)

	got, ok := c.Get("fresh")
	if !ok {
		t.Fatal("expected fresh entry to survive reap cycle")
	}
	if string(got) != "keep me" {
		t.Errorf("want %q, got %q", "keep me", string(got))
	}
}

func TestCache_StatsEmpty(t *testing.T) {
	c := NewCache(5 * time.Minute)
	stats := c.Stats()

	if stats.ItemCount != 0 {
		t.Errorf("expected 0 items, got %d", stats.ItemCount)
	}
	if stats.AverageLifetime != 0 {
		t.Errorf("expected 0 average lifetime, got %s", stats.AverageLifetime)
	}
}

func TestCache_StatsNonEmpty(t *testing.T) {
	c := NewCache(5 * time.Minute)
	c.Add("k1", []byte("v1"))

	time.Sleep(20 * time.Millisecond)
	stats := c.Stats()

	if stats.ItemCount != 1 {
		t.Fatalf("expected 1 item, got %d", stats.ItemCount)
	}
	if stats.AverageLifetime <= 0 {
		t.Errorf("expected positive average lifetime, got %s", stats.AverageLifetime)
	}
}
