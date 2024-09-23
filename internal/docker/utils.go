package docker

import (
	"github.com/docker/docker/api/types"
)

func FilterOnlyContainersInProject(projectName string, containers []types.Container) []types.Container {
	var filteredContainers []types.Container
	for _, container := range containers {
		if containerProject, ok := container.Labels["com.docker.compose.project"]; ok && containerProject == projectName {
			filteredContainers = append(filteredContainers, types.Container{
				ID:     container.ID[:10],
				Image:  container.Image,
				State:  container.State,
				Labels: container.Labels,
			})
		}
	}
	return filteredContainers
}

func GroupContainersByServiceName(filteredContainers []types.Container) map[string][]types.Container {
	groupedContainers := make(map[string][]types.Container)
	for _, container := range filteredContainers {
		serviceName := container.Labels["com.docker.compose.service"]
		groupedContainers[serviceName] = append(groupedContainers[serviceName], container)
	}
	return groupedContainers
}
