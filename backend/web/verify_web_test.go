package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStaticUnknownPath404(t *testing.T) {
	for _, path := range []string{"/index.html", "/favicon.ico", "/secret"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		Handler(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("%s: expected 404, got %d", path, rec.Code)
		}
	}
}

func TestStaticServesAppJSWithJSType(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	Handler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "javascript") {
		t.Fatalf("app.js content-type = %q, want javascript", ct)
	}
}

func TestStaticMethodNotAllowed(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	Handler(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
