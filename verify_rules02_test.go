package main

import (
	"testing"
)

func TestGroupThreeLabelsDetached(t *testing.T) {
	r1 := opsRule0301()
	_ = opsRule0302()
	for _, l := range r1.RequiredLabels {
		if l == "reviewed" {
			t.Fatalf("rule 0301 labels polluted: %v", r1.RequiredLabels)
		}
	}
}

func TestGroupFourLabelsDetached(t *testing.T) {
	r1 := opsRule0401()
	_ = opsRule0402()
	for _, l := range r1.RequiredLabels {
		if l == "reviewed" {
			t.Fatalf("rule 0401 labels polluted: %v", r1.RequiredLabels)
		}
	}
}

func TestGroupAggregationDetached(t *testing.T) {
	if len(opsRules04()) != 8 {
		t.Fatalf("group 04 has %d rules, want 8", len(opsRules04()))
	}
}

func TestAllRuleCodesPresent(t *testing.T) {
	rules := opsRules()
	seen := map[string]bool{}
	for _, r := range rules {
		if seen[r.Code] {
			t.Fatalf("duplicate rule code %q", r.Code)
		}
		seen[r.Code] = true
	}
	if len(rules) != 112 {
		t.Fatalf("rules count = %d, want 112", len(rules))
	}
}

func TestGroupThreePairTwoDetached(t *testing.T) {
	r1 := opsRule0305()
	_ = opsRule0306()
	for _, l := range r1.RequiredLabels {
		if l == "reviewed" {
			t.Fatalf("rule 0305 labels polluted: %v", r1.RequiredLabels)
		}
	}
}
