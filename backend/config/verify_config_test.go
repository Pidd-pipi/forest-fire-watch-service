package config

import (
	"testing"
)

func TestConfigLoadFallsBackOnInvalidPort(t *testing.T) {
	t.Setenv("PORT", "not-a-number")
	if c := Load(); c.Port != 8080 {
		t.Fatalf("expected fallback port 8080, got %d", c.Port)
	}
	t.Setenv("PORT", "")
	if c := Load(); c.Port != 8080 {
		t.Fatalf("empty PORT should fall back, got %d", c.Port)
	}
}

func TestConfigAddressFallsBack(t *testing.T) {
	if addr := (Config{Port: 0}).Address(); addr != ":8080" {
		t.Fatalf("zero port address = %q, want :8080", addr)
	}
}
