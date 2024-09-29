package main

import (
	"github.com/docker/docker/client"
	"github.com/r3labs/sse/v2"
	"log"
	"net/http"
	"os"
)

func main() {
	// Create a Docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Error creating Docker client: %s", err)
	}

	// Serve static files (CSS, JS, etc.)
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route to handle dynamic pages
	port := os.Getenv("PORT")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		RenderDocumentation(w, r, dockerClient)
	})

	// Set up SSE server
	sseServer := sse.New()
	sseServer.CreateStream("status-updates")
	http.HandleFunc("/events", sseServer.ServeHTTP)
	go PublishStatusUpdates(dockerClient, sseServer)

	// Start the server
	if port == "" {
		port = "8080"
	}
	log.Println("Serving Laebel documentation site at:")
	log.Println("")
	log.Println("  http://localhost:" + port + "/")
	log.Println("")
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println("INTERNAL SERVER ERROR: Could not start server.")
		log.Println("Cause:", err.Error())
		log.Println("Hint: Bind port " + port + " to another host port, or set the PORT environment variable to change port.")
	}
}
