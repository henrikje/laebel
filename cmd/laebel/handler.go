package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"net/http"
	"os"
)

var projectName string

func RenderDocumentation(w http.ResponseWriter, r *http.Request, dockerClient *client.Client) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Determine project to document
	if projectName == "" {
		// Get the current container ID
		containerID, err := GetContainerID()
		if err != nil {
			reportInternalServerError(w, err, "Could not determine current container ID", "Are you sure you are running Laebel as a container?")
			return
		}

		// Check if it’s part of a Compose project
		newProjectName, err := IsPartOfComposeProject(containerID, dockerClient)
		if err != nil {
			reportInternalServerError(w, err, "Could not determine current project name", "Ensure Laebel has the Docker socket mounted as a volume: \"/var/run/docker.sock:/var/run/docker.sock:ro\"")
			return
		}
		if newProjectName == "" {
			newProjectName = os.Getenv("COMPOSE_PROJECT_NAME")
		}
		if newProjectName == "" {
			reportNoProjectError(w)
			return
		}
		log.Println("Current project name:", newProjectName)
		projectName = newProjectName
	}

	// Get all containers
	containers, err := GetAllContainers(dockerClient)
	if err != nil {
		reportInternalServerError(w, err, "Unable to list containers", "")
		return
	}

	// Filter containers that are part of the project
	projectContainers := filterOnlyContainersInProject(containers, projectName)

	// Inspect each remaining container
	var projectContainersWithDetails []types.ContainerJSON
	for _, container := range projectContainers {
		containerDetails, err := InspectContainer(container.ID, dockerClient)
		if err != nil {
			reportInternalServerError(w, err, "Unable to inspect container", "")
			return
		}
		projectContainersWithDetails = append(projectContainersWithDetails, containerDetails)
	}

	// Transform containers to project
	project := TransformContainersToProject(projectContainersWithDetails, projectName)

	RenderDocument(w, err, project)
}

func filterOnlyContainersInProject(containers []types.Container, projectName string) []types.Container {
	var filteredContainers []types.Container
	for _, container := range containers {
		containerProject := container.Labels["com.docker.compose.project"]
		if containerProject == projectName {
			filteredContainers = append(filteredContainers, container)
		}
	}
	return filteredContainers
}

func reportInternalServerError(w http.ResponseWriter, err error, message string, hint string) {
	log.Println("INTERNAL SERVER ERROR:", message)
	log.Println("Cause:", err.Error())
	if hint != "" {
		log.Println("Hint:", hint)
	}
	http.Error(w, "INTERNAL SERVER ERROR: "+message+"\nCause: "+err.Error(), http.StatusInternalServerError)
	if hint != "" {
		_, _ = w.Write([]byte("Hint: " + hint))
	}
}

func reportNoProjectError(w http.ResponseWriter) {
	log.Println("BAD REQUEST: Current container is not part of a Docker Compose project.")
	log.Println("Hint: Add Laebel as a service in your Docker Compose project.")
	log.Println("Hint: If you want to run Laebel as a stand-alone container, specify the COMPOSE_PROJECT_NAME environment variable.")
	http.Error(w, "BAD REQUEST: Current container is not part of a Docker Compose project.", http.StatusBadRequest)
	_, _ = w.Write([]byte("Hint: Add Laebel as a service in your Docker Compose project.\n"))
	_, _ = w.Write([]byte("Hint: If you want to run Laebel as a stand-alone container, specify the COMPOSE_PROJECT_NAME environment variable."))
}
