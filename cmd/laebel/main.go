package main

import (
	"github.com/docker/docker/client"
	"github.com/r3labs/sse/v2"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Set up dependencies
	dockerClient := createDockerClient()
	tmpl := loadTemplates()

	// Determine current project
	projectName := determineCurrentProjectName(dockerClient)

	// Register handlers
	registerStaticFileHandler()
	registerPageHandler(projectName, tmpl, dockerClient)
	registerEventPublisher(dockerClient)

	// Start the server
	startServer()
}

func determineCurrentProjectName(dockerClient *client.Client) string {
	// Get the current container ID
	containerID, err := GetContainerID()
	if err != nil {
		// TODO Extract helper function for logging and hinting
		log.Fatalf("Could not determine current container ID\n"+
			"Cause: %s\n"+
			"Hint: Are you sure you are running Laebel as a container?", err)
	}

	// Check if it’s part of a Compose project
	projectName, err := IsPartOfComposeProject(containerID, dockerClient)
	if err != nil {
		log.Fatalf("Could not determine current project name\n"+
			"Cause: %s\n"+
			"Hint: Ensure Laebel has the Docker socket mounted as a volume: \"/var/run/docker.sock:/var/run/docker.sock:ro\"", err)
	}
	if projectName == "" {
		projectName = os.Getenv("COMPOSE_PROJECT_NAME")
	}
	if projectName == "" {
		log.Fatal("BAD REQUEST: Current container is not part of a Docker Compose project.\n" +
			"Hint: Add Laebel as a service in your Docker Compose project.\n" +
			"Hint: If you want to run Laebel as a stand-alone container, specify the COMPOSE_PROJECT_NAME environment variable.")
	}
	log.Println("Current project name:", projectName)
	return projectName
}

func loadTemplates() *template.Template {
	var escapeReplacer = strings.NewReplacer(".", "#46;", "(", "#40;", ")", "#41;")
	tmpl, err := template.New("index.html").Funcs(
		template.FuncMap{
			"escape": func(raw string) string { return escapeReplacer.Replace(raw) },
		}).ParseFiles(
		filepath.Join("web", "templates", "index.html"),
		filepath.Join("web", "templates", "main.html"),
		filepath.Join("web", "templates", "serviceGraph.html"),
		filepath.Join("web", "templates", "service.html"),
		filepath.Join("web", "templates", "serviceStatus.html"),
		filepath.Join("web", "templates", "serviceStatusSummary.html"),
		filepath.Join("web", "templates", "clipboard.html"),
	)
	if err != nil {
		log.Fatalf("INTERNAL SERVER ERROR: Unable to load template\nCause: %s", err)
	}
	return tmpl
}

func createDockerClient() *client.Client {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("INTERNAL SERVER ERROR: Error creating Docker client\nCause: %s", err)
	}
	return dockerClient
}

func registerStaticFileHandler() {
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func registerPageHandler(projectName string, tmpl *template.Template, dockerClient *client.Client) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HandleRequest(w, r, projectName, tmpl, dockerClient)
	})
}

func registerEventPublisher(dockerClient *client.Client) {
	sseServer := sse.New()
	sseServer.Headers = map[string]string{
		"Content-Type": "text/event-stream; charset=utf-8",
	}
	sseServer.CreateStream("refresh")
	http.HandleFunc("/events", sseServer.ServeHTTP)
	go PublishStatusUpdates(dockerClient, sseServer)
}

func startServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Serving Laebel documentation site at:")
	log.Println("")
	log.Println("  http://localhost:" + port + "/")
	log.Println("")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("INTERNAL SERVER ERROR: Could not start server.\n"+
			"Cause: %s\n"+
			"Hint: Bind port %s to another host port, or set the PORT environment variable to change port.", err, port)
	}
}
