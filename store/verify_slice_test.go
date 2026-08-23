package store

import (
	"testing"
)

func TestStoreInstancesIsolated(t *testing.T) {
	a := New()
	b := New()
	a.items[0].Status = "contained"
	if b.items[0].Status == "contained" {
		t.Fatalf("stores share seed backing: %v", b.items[0].Status)
	}
}

func TestStoreListSnapshotDetached(t *testing.T) {
	s := New()
	items := s.List()
	items[0].Status = "contained"
	if s.List()[0].Status == "contained" {
		t.Fatalf("list returned internal slice: %v", s.List()[0].Status)
	}
}
