package main

import (
	"github.com/docker/docker/api/types"
	"net/http"
)

func RenderDocumentation(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Get the current container ID
	containerID, err := GetContainerID()
	if err != nil {
		InternalServerError(w, err, "Could not determine current container ID", "Are you sure you are running docodash as a container?")
		return
	}
	println("Docodash container ID:", containerID)

	// Check if it’s part of a Compose project
	projectName, err := IsPartOfComposeProject(containerID)
	if err != nil {
		InternalServerError(w, err, "Could not determine current project name", "Ensure docodash has the Docker socket mounted as a volume: \"/var/run/docker.sock:/var/run/docker.sock:ro\"")
		return
	}
	println("Current project name:", projectName)
	if projectName == "" {
		NoProjectError(w)
		return
	}

	// Get all containers
	containers, err := GetAllContainers()
	if err != nil {
		InternalServerError(w, err, "Unable to list containers", "")
		return
	}

	// Remove current container from the list
	var currentContainer types.Container
	for i, container := range containers {
		if container.ID == containerID {
			currentContainer = container
			containers = append(containers[:i], containers[i+1:]...)
			break
		}
	}

	// Filter containers that are part of the project
	projectContainers := FilterOnlyContainersInProject(containers, projectName)

	// Transform containers to project
	project := TransformContainersToProject(projectContainers, currentContainer, projectName)

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
	println("INTERNAL SERVER ERROR:", message, err)
	if hint != "" {
		println("Hint:", hint)
	}
	http.Error(w, "INTERNAL SERVER ERROR: "+message+"\n\nCause: "+err.Error(), http.StatusInternalServerError)
	if hint != "" {
		_, _ = w.Write([]byte("Hint: " + hint))
	}
}

func NoProjectError(w http.ResponseWriter) {
	println("BAD REQUEST: Current container is not part of a Docker Compose project.")
	println("Hint: Are you running docodash as a service in a Docker Compose project?")
	http.Error(w, "BAD REQUEST: Current container is not part of a Docker Compose project.\n", http.StatusBadRequest)
	_, _ = w.Write([]byte("Hint: Are you running docodash as a service in a Docker Compose project?"))
}
