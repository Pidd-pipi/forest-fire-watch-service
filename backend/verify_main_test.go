package main

import (
	"testing"
)

func TestStartupConfigFallsBack(t *testing.T) {
	t.Setenv("PORT", "not-a-number")
	if c := startupConfig(); c.Port != 8080 {
		t.Fatalf("startup config port = %d, want 8080", c.Port)
	}
}
