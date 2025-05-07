package main

import (
	"log"
	"net/http"

	"github.com/RobertoPaulino/web-crawler/cmd/web/handler"
)

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle all other routes
	http.HandleFunc("/", handler.Handler)

	// Start server
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
