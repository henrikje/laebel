package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"os"
	"regexp"
	"strings"
)

// GetContainerID retrieves the current container ID
func GetContainerID() (string, error) {
	// cat /proc/self/mountinfo | grep "/docker/containers/" | head -1 | awk '{print $4}' | sed 's/\/var\/lib\/docker\/containers\///g' | sed 's/\/resolv.conf//g'

	// Read /proc/self/mountinfo
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

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// IsInComposeCluster checks if the container is part of a Docker Compose cluster
func IsInComposeCluster(containerID string) (bool, string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false, "", err
	}

	container, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return false, "", err
	}

	projectName, ok := container.Config.Labels["com.docker.compose.project"]
	if !ok {
		return false, "", nil
	}

	return true, projectName, nil
}
