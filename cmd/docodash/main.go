package main

import (
	"log"
	"net/http"
)

func main() {
	// Serve static files (CSS, JS, etc.)
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route to handle dynamic pages
	http.HandleFunc("/", RenderDocumentation)

	// Start the server
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
