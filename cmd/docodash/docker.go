package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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

func GetAllContainers() ([]types.Container, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// IsPartOfComposeProject checks if the container is part of a Docker Compose cluster
func IsPartOfComposeProject(containerID string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	containerInfo, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return "", err
	}

	projectName, ok := containerInfo.Config.Labels["com.docker.compose.project"]
	if !ok {
		return "", nil
	}

	return projectName, nil
}
