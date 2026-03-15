package main

import (
	"log"
	"net/http"

	"example.com/example-project/internal/httpserver"
)

func main() {
	addr := ":8080"
	handler := httpserver.NewMux()

	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
