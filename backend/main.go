package main

import (
	"example.com/forest-fire-watch-service/config"
	"example.com/forest-fire-watch-service/health"
	"example.com/forest-fire-watch-service/httpapi"
	"example.com/forest-fire-watch-service/store"
	"example.com/forest-fire-watch-service/web"
	"log"
	"net/http"
)

func main() {
	c := config.Load()
	m := http.NewServeMux()
	m.HandleFunc("/healthz", health.Handler)
	m.Handle("/api/v1/", httpapi.New(store.New()))
	m.HandleFunc("/", web.Handler)
	log.Printf("forest fire watch listening on %s", c.Address())
	log.Fatal(serveAddress(c.Address(), m))
}
