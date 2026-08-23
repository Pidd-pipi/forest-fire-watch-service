package httpapi

import (
	"bytes"
	"example.com/forest-fire-watch-service/store"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPPostValidPathAccepted(t *testing.T) {
	h := New(store.New())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/alerts/fa-301/status", bytes.NewBufferString(`{"status":"contained"}`))
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("valid path rejected: %d", rec.Code)
	}
}

func TestHTTPUpdateStatusPersists(t *testing.T) {
	st := store.New()
	h := New(st)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/alerts/fa-301/status", bytes.NewBufferString(`{"status":"contained"}`))
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("update failed: %d", rec.Code)
	}
	found := false
	for _, it := range st.List() {
		if it.ID == "fa-301" && it.Status == "contained" {
			found = true
		}
	}
	if !found {
		t.Fatalf("status not persisted: %+v", st.List())
	}
}
