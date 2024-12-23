package main

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
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
	currentContainerID := detectCurrentContainerID()
	projectName := determineComposeProjectName(dockerClient, currentContainerID)

	// Register handlers
	registerStaticFileHandler()
	registerPageHandler(projectName, tmpl, dockerClient)
	registerEventPublisher(dockerClient)

	// Start the server
	containerPort, hostPort := detectPorts(dockerClient, currentContainerID)
	startServer(containerPort, hostPort, projectName)
}

func createDockerClient() *client.Client {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fatal(err, "Error creating Docker client")
	}
	return dockerClient
}

func loadTemplates() *template.Template {
	var escapeReplacer = strings.NewReplacer(".", "#46;", "(", "#40;", ")", "#41;", "|", "#124;", "[", "#91;", "]", "#93;", "{", "#123;", "}", "#125;", "`", "'")
	tmpl, err := template.New("index.html").Funcs(
		template.FuncMap{
			"escape": func(raw string) string { return escapeReplacer.Replace(raw) },
		}).ParseFiles(
		filepath.Join("web", "templates", "index.html"),
		filepath.Join("web", "templates", "resources.html"),
		filepath.Join("web", "templates", "graph.html"),
		filepath.Join("web", "templates", "graph.mmd.html"),
		filepath.Join("web", "templates", "service.html"),
		filepath.Join("web", "templates", "serviceStatus.html"),
		filepath.Join("web", "templates", "volume.html"),
		filepath.Join("web", "templates", "network.html"),
		filepath.Join("web", "templates", "clipboard.html"),
	)
	if err != nil {
		fatal(err, "Unable to load template")
	}
	return tmpl
}

func detectCurrentContainerID() string {
	// Get the current container ID
	containerID, err := GetContainerID()
	if err != nil {
		fatal(err, "Could not determine current container ID", "Are you sure you are running Laebel as a container?")
	}
	return containerID
}

func determineComposeProjectName(dockerClient *client.Client, containerID string) string {
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

func registerStaticFileHandler() {
	fileHandler := http.FileServer(http.Dir("web/static"))
	cachingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=31536000") // 1 year
		fileHandler.ServeHTTP(w, r)
	})
	http.Handle("/static/", http.StripPrefix("/static/", cachingHandler))
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
	stream := sseServer.CreateStream("updates")
	// AutoReplay is unnecessary as the client starts with up-to-date information.
	// We also don't want to replay old "reload" events.
	stream.AutoReplay = false
	http.HandleFunc("/events", sseServer.ServeHTTP)
	go PublishStatusUpdates(dockerClient, sseServer)
}

func detectPorts(dockerClient *client.Client, currentContainerID string) (string, string) {
	// Determine port exposed from container
	containerPort := os.Getenv("PORT")
	if containerPort == "" {
		containerPort = "8000"
	}

	// Look up the corresponding port exposed on host
	currentContainer, err := dockerClient.ContainerInspect(context.Background(), currentContainerID)
	if err != nil {
		fatal(err, "Could not inspect current container.", "Ensure Laebel has the Docker socket mounted as a volume: \"/var/run/docker.sock:/var/run/docker.sock:ro\"")
	}
	portKey := nat.Port(containerPort + "/tcp")
	portBindings := currentContainer.NetworkSettings.Ports[portKey]
	hostPort := ""
	if len(portBindings) > 0 {
		hostPort = portBindings[0].HostPort
	}
	return containerPort, hostPort
}

func startServer(containerPort, hostPort, projectName string) {
	if hostPort == "" {
		log.Println("Warning: No host port bound to container")
		log.Printf("Hint: Bind port " + containerPort + " to a port on the host, either using 'ports' in the Docker Compose file or with the -p flag.")
		hostPort = containerPort
	} else {
		log.Printf("Serving Laebel documentation site for project '%s' at:\n", projectName)
		log.Println("")
		log.Println("  http://localhost:" + hostPort + "/")
		log.Println("")
	}
	err := http.ListenAndServe(":"+containerPort, nil)
	if err != nil {
		fatal(err, "Could not start server", "Bind port "+containerPort+" to another host port than "+hostPort+", or set the PORT environment variable to change the exposed port.")
	}
}
