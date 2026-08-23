package httpapi

import (
	"bytes"
	"example.com/forest-fire-watch-service/store"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAlertAPI(t *testing.T) {
	s := httptest.NewServer(New(store.New()))
	defer s.Close()
	for _, tc := range []struct {
		method, path, body string
		want               int
	}{{"GET", "/api/v1/alerts", "", 200}, {"POST", "/api/v1/alerts/fa-301/status", `{"status":"contained"}`, 200}, {"POST", "/api/v1/alerts/fa-301/status", `{"status":"burning"}`, 400}, {"POST", "/api/v1/alerts/unknown/status", `{"status":"resolved"}`, 404}} {
		req, _ := http.NewRequest(tc.method, s.URL+tc.path, bytes.NewBufferString(tc.body))
		res, e := http.DefaultClient.Do(req)
		if e != nil || res.StatusCode != tc.want {
			t.Fatalf("%s %s: err=%v status=%d", tc.method, tc.path, e, res.StatusCode)
		}
		res.Body.Close()
	}
}
