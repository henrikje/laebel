package main

import (
	"github.com/docker/docker/api/types"
	"os"
)

func TransformContainersToProject(containers []types.Container, projectName string) Project {
	// TODO Where do I get the project description from?
	containersByServiceName := make(map[string][]types.Container)

	// Group containers by service name
	for _, container := range containers {
		serviceName := container.Labels["com.docker.compose.service"]
		containersByServiceName[serviceName] = append(containersByServiceName[serviceName], container)
	}

	// Create a service for each group of containers
	var services []Service
	for serviceName, serviceContainers := range containersByServiceName {
		container := serviceContainers[0] // Use the first container to extract metadata
		image := container.Image
		links := make([]Link, 0)
		url := container.Labels["org.opencontainers.image.url"]
		if url != "" {
			links = append(links, Link{Label: "Website", URL: url})
		}
		documentation := container.Labels["org.opencontainers.image.documentation"]
		if documentation != "" {
			links = append(links, Link{Label: "Documentation", URL: documentation})
		}
		source := container.Labels["org.opencontainers.image.source"]
		if source != "" {
			links = append(links, Link{Label: "Source code", URL: source})
		}
		containerStructs := make([]Container, 0)
		for _, serviceContainer := range serviceContainers {
			containerStructs = append(containerStructs, Container{
				ID:     serviceContainer.ID,
				Names:  serviceContainer.Names,
				Status: serviceContainer.Status,
			})
		}
		status := Status{
			Created:    0,
			Running:    0,
			Paused:     0,
			Restarting: 0,
			Exited:     0,
			Removing:   0,
			Dead:       0,
			Stopped:    0,
		}
		for index, _ := range containerStructs {
			switch serviceContainers[index].State {
			case "created":
				status.Created++
			case "running":
				status.Running++
			case "paused":
				status.Paused++
			case "restarting":
				status.Restarting++
			case "exited":
				status.Exited++
			case "removing":
				status.Removing++
			case "dead":
				status.Dead++
			case "stopped":
				status.Stopped++
			}
		}
		service := Service{
			Name:        serviceName,
			Title:       container.Labels["org.opencontainers.image.projectTitle"],
			Description: container.Labels["org.opencontainers.image.description"],
			Image:       image,
			Status:      status,
			Links:       links,
			Containers:  containerStructs,
		}
		services = append(services, service)
	}
	return Project{
		Name:        projectName,
		Title:       os.Getenv("DOCODASH_PROJECT_TITLE"),
		Description: os.Getenv("DOCODASH_PROJECT_DESCRIPTION"),
		Services:    services,
	}
}
