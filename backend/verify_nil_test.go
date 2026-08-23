package main

import (
	"testing"
)

func TestOpsQueryDefaultsReasonable(t *testing.T) {
	q := opsQueryDefaults(OpsQuery{})
	if q.Page != 1 {
		t.Fatalf("page not defaulted: %d", q.Page)
	}
	if q.PageSize < 1 || q.PageSize > 200 {
		t.Fatalf("page size not defaulted: %d", q.PageSize)
	}
}

func TestOpsBoundsClampsBeyondEnd(t *testing.T) {
	start, end := opsBounds(10, 99, 25)
	if start > 10 || start < 0 || end < start || end > 10 {
		t.Fatalf("bounds out of range: start=%d end=%d", start, end)
	}
	start2, end2 := opsBounds(0, 3, 25)
	if start2 != 0 || end2 != 0 {
		t.Fatalf("empty total should yield empty bounds: %d %d", start2, end2)
	}
}

func TestOpsClonePageIndependentSnapshot(t *testing.T) {
	src := OpsPage{Items: []OpsRecord{{ID: "r1", Status: OpsStatusQueued}}}
	cloned := opsClonePage(src)
	cloned.Items[0].Status = OpsStatusClosed
	if src.Items[0].Status != OpsStatusQueued {
		t.Fatalf("clone shares backing slice: %v", src.Items[0].Status)
	}
}

func TestOpsStatusValidRejectsEmpty(t *testing.T) {
	if opsStatusValid("") {
		t.Fatal("empty status should be invalid")
	}
}
