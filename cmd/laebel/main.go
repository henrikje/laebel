package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Serve static files (CSS, JS, etc.)
	fs := http.FileServer(http.Dir("web/static"))
	// TODO How do we serve static files with the correct MIME type?
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route to handle dynamic pages
	port := os.Getenv("PORT")
	http.HandleFunc("/", RenderDocumentation)

	// Start the server
	log.Println("Starting server on :" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
