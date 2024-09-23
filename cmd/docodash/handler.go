package main

import (
	"docodash/internal/docker"
	"github.com/docker/docker/api/types"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

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
	if err != nil {
		http.Error(w, "Failed to determine Compose project: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !inComposeCluster {
		http.Error(w, "Not part of a Docker Compose cluster!", http.StatusBadRequest)
		return
	}

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
	containersGroupedByServiceName := docker.GroupContainersByServiceName(filteredContainers)
	println("containersGroupedByServiceName:", containersGroupedByServiceName)

	// Extract links from container labels
	linksByService := ExtractServiceLinksFromLabels(containersGroupedByServiceName)
	println("linksByService:", linksByService)

	// Render template
	tmpl, err := template.ParseFiles(filepath.Join("web", "templates", "index.html"))
	if err != nil {
		http.Error(w, "Unable to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := struct {
		ContainersGroupedByServiceName map[string][]types.Container
		LinksByService                 map[string][]map[string]string
	}{
		ContainersGroupedByServiceName: containersGroupedByServiceName,
		LinksByService:                 linksByService,
	}

	err = tmpl.Execute(w, templateData)
	if err != nil {
		http.Error(w, "Unable to render template: "+err.Error(), http.StatusInternalServerError)
	}
}

func ExtractServiceLinksFromLabels(containersGroupedByServiceName map[string][]types.Container) map[string][]map[string]string {
	// Links are defined in container labels as follows:
	// net.henko.docodash.link.<key>.url = <url>
	// net.henko.docodash.link.<key>.label = <label>
	linksByService := make(map[string][]map[string]string)
	for serviceName, containers := range containersGroupedByServiceName {
		container := containers[0] // We only need one container to extract the links - they should all be the same
		for key, value := range container.Labels {
			if !strings.HasPrefix(key, "net.henko.docodash.link.") {
				continue
			}

			parts := strings.Split(key, ".")
			if len(parts) != 6 {
				println("Invalid link:", key)
				continue
			}

			linkKey := parts[4]
			linkType := parts[5]

			if linkType == "url" {
				if _, ok := linksByService[serviceName]; !ok {
					linksByService[serviceName] = make([]map[string]string, 0)
				}
				label := container.Labels["net.henko.docodash.link."+linkKey+".label"]
				println("link", linkKey, "url:", value, "label:", label)
				link := map[string]string{
					"url":   value,
					"label": label,
				}
				linksByService[serviceName] = append(linksByService[serviceName], link)
			}
		}
	}
	return linksByService
}
