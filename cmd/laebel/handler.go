package main

import (
	"github.com/docker/docker/api/types"
	"log"
	"net/http"
)

var containerID string
var projectName string

func RenderDocumentation(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Get the current container ID
	if containerID == "" {
		newContainerID, err := GetContainerID()
		if err != nil {
			InternalServerError(w, err, "Could not determine current container ID", "Are you sure you are running Laebel as a container?")
			return
		}
		log.Println("Laebel container ID:", newContainerID)
		containerID = newContainerID
	}

	// Check if it’s part of a Compose project
	if projectName == "" {
		newProjectName, err := IsPartOfComposeProject(containerID)
		if err != nil {
			InternalServerError(w, err, "Could not determine current project name", "Ensure Laebel has the Docker socket mounted as a volume: \"/var/run/docker.sock:/var/run/docker.sock:ro\"")
			return
		}
		if newProjectName == "" {
			NoProjectError(w)
			return
		}
		log.Println("Current project name:", newProjectName)
		projectName = newProjectName
	}

	// Get all containers
	containers, err := GetAllContainers()
	if err != nil {
		InternalServerError(w, err, "Unable to list containers", "")
		return
	}

	// Filter containers that are part of the project
	projectContainers := FilterOnlyContainersInProject(containers, projectName)

	// Transform containers to project
	project := TransformContainersToProject(projectContainers, projectName)

	RenderDocument(w, err, project)
}

func FilterOnlyContainersInProject(containers []types.Container, projectName string) []types.Container {
	var filteredContainers []types.Container
	for _, container := range containers {
		containerProject := container.Labels["com.docker.compose.project"]
		if containerProject == projectName {
			filteredContainers = append(filteredContainers, container)
		}
	}
	return filteredContainers
}

func InternalServerError(w http.ResponseWriter, err error, message string, hint string) {
	log.Println("INTERNAL SERVER ERROR:", message+":", err)
	if hint != "" {
		log.Println("Hint:", hint)
	}
	http.Error(w, "INTERNAL SERVER ERROR: "+message+"\n\nCause: "+err.Error(), http.StatusInternalServerError)
	if hint != "" {
		_, _ = w.Write([]byte("Hint: " + hint))
	}
}

func NoProjectError(w http.ResponseWriter) {
	log.Println("BAD REQUEST: Current container is not part of a Docker Compose project.")
	log.Println("Hint: Are you running Laebel as a service in a Docker Compose project?")
	http.Error(w, "BAD REQUEST: Current container is not part of a Docker Compose project.\n", http.StatusBadRequest)
	_, _ = w.Write([]byte("Hint: Are you running Laebel as a service in a Docker Compose project?"))
}
