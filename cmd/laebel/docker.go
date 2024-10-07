package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"os"
	"regexp"
	"strings"
)

func GetContainerID() (string, error) {

	data, err := os.ReadFile("/proc/self/mountinfo")
	if err != nil {
		return "", err
	}

	// Extract id from a line containing "/docker/containers/"
	for _, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, "/docker/containers/") {
			regex := regexp.MustCompile("/docker/containers/(?P<id>[^/]+)/")
			return regex.FindStringSubmatch(line)[1], nil
		}
	}

	return "", nil
}

func GetAllContainersInProject(projectName string, dockerClient *client.Client) ([]types.Container, error) {
	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.KeyValuePair{Key: "label", Value: "com.docker.compose.project=" + projectName},
		),
	})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func InspectEachContainer(containers []types.Container, dockerClient *client.Client) ([]types.ContainerJSON, error) {
	var containerInfos []types.ContainerJSON
	for _, theContainer := range containers {
		containerInfo, err := dockerClient.ContainerInspect(context.Background(), theContainer.ID)
		if err != nil {
			return nil, err
		}
		containerInfos = append(containerInfos, containerInfo)
	}
	return containerInfos, nil
}

// IsPartOfComposeProject checks if the container is part of a Docker Compose cluster
func IsPartOfComposeProject(containerID string, dockerClient *client.Client) (string, error) {
	containerInfo, err := dockerClient.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return "", err
	}

	projectName, ok := containerInfo.Config.Labels["com.docker.compose.project"]
	if !ok {
		return "", nil
	}

	return projectName, nil
}
