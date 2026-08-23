package main

import (
	"context"
	"testing"
	"time"
)

func TestOpsContextKeepsParentDeadline(t *testing.T) {
	parent, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	ctx, cancel2 := opsContext(parent, 5*time.Second)
	defer cancel2()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("no deadline propagated")
	}
	if time.Until(deadline) > time.Second {
		t.Fatalf("parent deadline lost: %v", time.Until(deadline))
	}
}

func TestOpsDelayHonorsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	start := time.Now()
	err := opsDelay(ctx, 2*time.Second)
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if time.Since(start) > time.Second {
		t.Fatalf("delay not interrupted by cancel: %v", time.Since(start))
	}
}

func TestOpsTransitionHonorsParentContext(t *testing.T) {
	s := newOpsService([]OpsRecord{{ID: "r1", Subject: "s", Owner: "o", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "s1"}}})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := s.Transition(ctx, "r1", 1, OpsStatusActive, "tester"); err == nil {
		t.Fatal("transition should fail when parent context is cancelled")
	}
}

func TestOpsCreateHonorsParentContext(t *testing.T) {
	s := newOpsService(nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	rec := OpsRecord{ID: "r1", Subject: "s", Owner: "o", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "s1"}}
	if _, err := s.Create(ctx, rec); err == nil {
		t.Fatal("create should fail when parent context is cancelled")
	}
}
