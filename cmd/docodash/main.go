package main

import (
	"docodash/internal/docker"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Get the current container ID
	containerID, err := docker.GetContainerID()
	println("containerID:", containerID)
	if err != nil {
		http.Error(w, "Unable to determine current container ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if it’s part of a Compose cluster
	inComposeCluster, projectName, err := docker.IsInComposeCluster(containerID)
	println("inComposeCluster:", inComposeCluster)
	println("projectName:", projectName)
	// if err != nil || !inComposeCluster {
	// 	http.Error(w, "Not part of a Docker Compose cluster", http.StatusBadRequest)
	// 	return
	// }

	// Get all containers
	containers, err := docker.GetAllContainers()
	println("containers:", containers)
	if err != nil {
		http.Error(w, "Unable to list containers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter containers by the same Compose project
	filteredContainers := docker.FilterOnlyContainersInProject(projectName, containers)
	println("filteredContainers:", filteredContainers)

	// Group containers by service name
	containersGroupedByName := docker.GroupContainersByServiceName(filteredContainers)
	println("containersGroupedByName:", containersGroupedByName)

	// Render template
	tmpl, err := template.ParseFiles(filepath.Join("web", "templates", "index.html"))
	if err != nil {
		http.Error(w, "Unable to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, containersGroupedByName)
	if err != nil {
		http.Error(w, "Unable to render template: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Serve static files (CSS, JS, etc.)
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route to handle dynamic pages
	http.HandleFunc("/", handler)

	// Start the server
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
