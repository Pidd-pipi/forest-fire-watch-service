package main

import (
	"testing"
)

func TestOpsRulesTotalCountIs112(t *testing.T) {
	rules := opsRules()
	if len(rules) != 112 {
		t.Fatalf("rules count = %d, want 112", len(rules))
	}
}

func TestOpsRules01Complete(t *testing.T) {
	if len(opsRules01()) != 8 {
		t.Fatalf("group 01 has %d rules, want 8", len(opsRules01()))
	}
}

func TestOpsRule0101LabelsIsolated(t *testing.T) {
	r1 := opsRule0101()
	_ = opsRule0102()
	for _, l := range r1.RequiredLabels {
		if l == "reviewed" {
			t.Fatalf("rule 0101 labels polluted: %v", r1.RequiredLabels)
		}
	}
}

func TestOpsRule0201LabelsIsolated(t *testing.T) {
	r1 := opsRule0201()
	_ = opsRule0202()
	for _, l := range r1.RequiredLabels {
		if l == "reviewed" {
			t.Fatalf("rule 0201 labels polluted: %v", r1.RequiredLabels)
		}
	}
}

func TestOpsRule0203LabelsIsolated(t *testing.T) {
	r1 := opsRule0203()
	_ = opsRule0204()
	for _, l := range r1.RequiredLabels {
		if l == "reviewed" {
			t.Fatalf("rule 0203 labels polluted: %v", r1.RequiredLabels)
		}
	}
}
