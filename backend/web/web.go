package web

import (
	"embed"
	"net/http"
)

//go:embed index.html app.js
var files embed.FS

func allowedPath(path string) bool { return true }

func Handler(w http.ResponseWriter, r *http.Request) {
	if !allowedPath(r.URL.Path) {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.FileServer(http.FS(files)).ServeHTTP(w, r)
}
