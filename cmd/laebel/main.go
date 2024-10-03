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
	startServer(projectName)
}

func determineCurrentProjectName(dockerClient *client.Client) string {
	// Get the current container ID
	containerID, err := GetContainerID()
	if err != nil {
		fatal(err, "Could not determine current container ID", "Are you sure you are running Laebel as a container?")
	}

	// Check if it’s part of a Compose project
	projectName, err := IsPartOfComposeProject(containerID, dockerClient)
	if err != nil {
		fatal(err, "Could not determine Docker Compose project name", "Ensure Laebel has the Docker socket mounted as a volume: \"/var/run/docker.sock:/var/run/docker.sock:ro\"")
	}
	if projectName == "" {
		projectName = os.Getenv("COMPOSE_PROJECT_NAME")
	}
	if projectName == "" {
		fatal(nil, "No Docker Compose project detected.", "Add Laebel as a service in your Docker Compose project.", "If you want to run Laebel as a stand-alone container, specify the COMPOSE_PROJECT_NAME environment variable.")
	}
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
		fatal(err, "Unable to load template")
	}
	return tmpl
}

func createDockerClient() *client.Client {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fatal(err, "Error creating Docker client")
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
	sseServer.CreateStream("updates")
	http.HandleFunc("/events", sseServer.ServeHTTP)
	go PublishStatusUpdates(dockerClient, sseServer)
}

func startServer(projectName string) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Serving Laebel documentation site for project '%s' at:\n", projectName)
	log.Println("")
	log.Println("  http://localhost:" + port + "/")
	log.Println("")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fatal(err, "Could not start server", "Bind port "+port+" to another host port, or set the PORT environment variable to change port.")
	}
}
