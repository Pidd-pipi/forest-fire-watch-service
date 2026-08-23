package health

import (
	"encoding/json"
	"net/http"
)

func methodAllowed(r *http.Request) bool { return true }

func Handler(w http.ResponseWriter, r *http.Request) {
	if !methodAllowed(r) {
		http.Error(w, "method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
