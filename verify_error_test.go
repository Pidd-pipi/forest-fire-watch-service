package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrapOpsPreservesOpsErrorType(t *testing.T) {
	err := wrapOps("store.put", "create", ErrOpsNotFound)
	if !errors.Is(err, ErrOpsNotFound) {
		t.Fatalf("wrapOps broke sentinel chain: %v", err)
	}
}

func TestOpsCodeClassifiesWrappedErrors(t *testing.T) {
	if code := opsCode(fmt.Errorf("outer: %w", ErrOpsConflict)); code != "conflict" {
		t.Fatalf("expected conflict, got %s", code)
	}
	if code := opsCode(fmt.Errorf("outer: %w", ErrOpsNotFound)); code != "not_found" {
		t.Fatalf("expected not_found, got %s", code)
	}
}

func TestOpsIsNotFoundMatchesWrappedError(t *testing.T) {
	if !opsIsNotFound(fmt.Errorf("layer: %w", ErrOpsNotFound)) {
		t.Fatalf("wrapped not found not recognized")
	}
}

func TestOpsIsConflictMatchesWrappedError(t *testing.T) {
	if !opsIsConflict(fmt.Errorf("layer: %w", ErrOpsConflict)) {
		t.Fatalf("wrapped conflict not recognized")
	}
}

func TestEnterpriseMiddlewareLatencyHeader(t *testing.T) {
	srv := httptest.NewServer(opsEnterpriseMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.Header.Get("X-Operations-Latency-Ms") == "" {
		t.Fatalf("latency header missing on real response")
	}
}
