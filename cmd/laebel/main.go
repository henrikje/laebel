package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Serve static files (CSS, JS, etc.)
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route to handle dynamic pages
	port := os.Getenv("PORT")
	http.HandleFunc("/", RenderDocumentation)

	// Start the server
	if port == "" {
		port = "8080"
	}
	log.Println("Serving Laebel documentation site at:")
	log.Println("")
	log.Println("  http://localhost:" + port + "/")
	log.Println("")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println("INTERNAL SERVER ERROR: Could not start server.")
		log.Println("Cause:", err.Error())
		log.Println("Hint: Set the PORT environment variable to change port.")
	}
}
