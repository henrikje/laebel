package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func TransformContainersToProject(projectContainers []types.ContainerJSON, projectName string) Project {
	containers := filterOutHiddenContainers(projectContainers)
	containersByServiceName := groupContainersByServiceName(containers)
	services := createServiceForEachServiceName(containersByServiceName)
	serviceGroups := groupServicesByGroup(containers, services)
	projectLinks := extractProjectLinks()
	return Project{
		Name:          projectName,
		Title:         os.Getenv("LAEBEL_PROJECT_TITLE"),
		Description:   os.Getenv("LAEBEL_PROJECT_DESCRIPTION"),
		Links:         projectLinks,
		Icon:          os.Getenv("LAEBEL_PROJECT_ICON"),
		ServiceGroups: serviceGroups,
	}
}

func filterOutHiddenContainers(projectContainers []types.ContainerJSON) []types.ContainerJSON {
	noHiddenContainers := make([]types.ContainerJSON, 0)
	for _, container := range projectContainers {
		if container.Config.Labels["net.henko.laebel.hidden"] != "true" {
			noHiddenContainers = append(noHiddenContainers, container)
		}
	}
	return noHiddenContainers
}

func groupContainersByServiceName(containers []types.ContainerJSON) map[string][]types.ContainerJSON {
	containersByServiceName := make(map[string][]types.ContainerJSON)
	for _, container := range containers {
		serviceName := container.Config.Labels["com.docker.compose.service"]
		containersByServiceName[serviceName] = append(containersByServiceName[serviceName], container)
	}
	return containersByServiceName
}

func createServiceForEachServiceName(containersByServiceName map[string][]types.ContainerJSON) []Service {
	var services []Service
	for serviceName, serviceContainers := range containersByServiceName {
		service := transformContainersToService(serviceContainers, serviceName)
		services = append(services, service)
	}
	return services
}

func groupServicesByGroup(containers []types.ContainerJSON, services []Service) []ServiceGroup {
	// Extract group-to-service mapping
	groupNameByServiceName := make(map[string]string)
	for _, container := range containers {
		groupName := container.Config.Labels["net.henko.laebel.group"]
		serviceName := container.Config.Labels["com.docker.compose.service"]
		groupNameByServiceName[serviceName] = groupName
	}

	// Group services by category
	servicesByGroupName := make(map[string][]Service)
	for _, service := range services {
		category := groupNameByServiceName[service.Name]
		servicesByGroupName[category] = append(servicesByGroupName[category], service)
	}

	// Create a service group for each group of services
	var serviceGroups []ServiceGroup
	for category, categoryServices := range servicesByGroupName {
		sort.Slice(categoryServices, func(i, j int) bool {
			return categoryServices[i].Name < categoryServices[j].Name
		})
		serviceGroup := ServiceGroup{
			Name:     category,
			Services: categoryServices,
		}
		serviceGroups = append(serviceGroups, serviceGroup)
	}

	// Sort service groups by name
	sort.Slice(serviceGroups, func(i, j int) bool {
		return serviceGroups[i].Name < serviceGroups[j].Name
	})

	// Special treatment for nameless categories
	if len(serviceGroups) == 1 && serviceGroups[0].Name == "" {
		serviceGroups[0].Name = "Services"
	} else {
		for i := range serviceGroups {
			if serviceGroups[i].Name == "" {
				serviceGroups[i].Name = "Services"
			}
		}
	}
	return serviceGroups
}

func extractProjectLinks() []Link {
	projectLinks := make([]Link, 0)
	url := os.Getenv("LAEBEL_PROJECT_URL")
	if url != "" {
		projectLinks = append(projectLinks, Link{Title: "Website", URL: url})
	}
	documentation := os.Getenv("LAEBEL_PROJECT_DOCUMENTATION")
	if documentation != "" {
		projectLinks = append(projectLinks, Link{Title: "Documentation", URL: documentation})
	}
	source := os.Getenv("LAEBEL_PROJECT_SOURCE")
	if source != "" {
		projectLinks = append(projectLinks, Link{Title: "Source code", URL: source})
	}
	return projectLinks
}

func transformContainersToService(serviceContainers []types.ContainerJSON, serviceName string) Service {
	container := serviceContainers[0] // Use the first container to extract metadata
	links := extractServiceLinks(container)
	containerStructs := make([]Container, 0)
	for _, serviceContainer := range serviceContainers {
		containerHealth := "unknown"
		if serviceContainer.State.Health != nil {
			health := *serviceContainer.State.Health
			containerHealth = health.Status
		}
		parsedDate, err := time.Parse(time.RFC3339Nano, serviceContainer.Created)
		var created string
		if err != nil {
			created = serviceContainer.Created
		} else {
			created = parsedDate.Format(time.DateTime)
		}
		containerStructs = append(containerStructs, Container{
			ID:      serviceContainer.ID[:12],
			Name:    strings.TrimLeft(serviceContainer.Name, "/"),
			Created: created,
			Status:  serviceContainer.State.Status,
			Health:  containerHealth,
			Ports:   extractContainerPorts(serviceContainer.HostConfig.PortBindings),
		})
		sort.Slice(containerStructs, func(i, j int) bool {
			return containerStructs[i].Name < containerStructs[j].Name
		})
	}
	status := extractStatus(containerStructs, serviceContainers)
	dependsOn := extractDependsOn(container)
	service := Service{
		Name:        serviceName,
		Title:       container.Config.Labels["org.opencontainers.image.title"],
		Description: container.Config.Labels["org.opencontainers.image.description"],
		Image:       container.Config.Image,
		Status:      status,
		Links:       links,
		Ports:       extractServicePorts(serviceContainers),
		DependsOn:   dependsOn,
		Containers:  containerStructs,
	}
	return service
}

func extractServicePorts(containers []types.ContainerJSON) []Port {
	ports := make([]Port, 0)
	for _, container := range containers {
		for _, portBindings := range container.HostConfig.PortBindings {
			for _, portBinding := range portBindings {
				port := Port{
					Number:      portBinding.HostPort,
					Description: container.Config.Labels["net.henko.laebel.port."+portBinding.HostPort+".description"],
				}
				ports = append(ports, port)
			}
		}
	}
	sort.Slice(ports, func(i, j int) bool {
		inum, err := strconv.Atoi(ports[i].Number)
		if err != nil {
			return false
		}
		jnum, err := strconv.Atoi(ports[j].Number)
		if err != nil {
			return true
		}
		return inum < jnum
	})
	return ports
}

func extractContainerPorts(portMap nat.PortMap) []string {
	ports := make([]string, 0)
	for port, portBinding := range portMap {
		for _, hostPort := range portBinding {
			var hostString string
			if hostPort.HostIP == "" {
				hostString = ""
			} else {
				hostString = hostPort.HostIP + ":"
			}
			var portString string
			if hostPort.HostPort == port.Port() {
				portString = string(port)
			} else {
				portString = hostPort.HostPort + "->" + string(port)
			}
			portString = strings.TrimSuffix(portString, "/tcp")
			ports = append(ports, hostString+portString)
		}
	}
	return ports
}

func extractServiceLinks(container types.ContainerJSON) []Link {
	links := make([]Link, 0)
	// Extract OpenContainers links
	url := container.Config.Labels["org.opencontainers.image.url"]
	if url != "" {
		links = append(links, Link{Title: "Website", URL: url})
	}
	documentation := container.Config.Labels["org.opencontainers.image.documentation"]
	if documentation != "" {
		links = append(links, Link{Title: "Documentation", URL: documentation})
	}
	source := container.Config.Labels["org.opencontainers.image.source"]
	if source != "" {
		links = append(links, Link{Title: "Source code", URL: source})
	}
	// Extract laebel links
	// net.henko.laebel.link.<key>.url
	// net.henko.laebel.link.<key>.title
	for key, value := range container.Config.Labels {
		if len(key) > 22 && key[:22] == "net.henko.laebel.link." {
			linkKey := key[22:]
			if linkKey[len(linkKey)-4:] == ".url" {
				titleKey := "net.henko.laebel.link." + linkKey[:len(linkKey)-4] + ".title"
				title := container.Config.Labels[titleKey]
				if title == "" {
					title = linkKey[:len(linkKey)-4]
				}
				links = append(links, Link{Title: title, URL: value})
			}
		}
	}
	return links
}

func extractStatus(containerStructs []Container, serviceContainers []types.ContainerJSON) Status {
	status := Status{
		Created:          0,
		Running:          0,
		RunningHealthy:   0,
		RunningUnhealthy: 0,
		Paused:           0,
		Restarting:       0,
		Exited:           0,
		Removing:         0,
		Dead:             0,
	}
	for index := range containerStructs {
		container := serviceContainers[index]
		switch container.State.Status {
		// Can be one of "created", "running", "paused", "restarting", "removing", "exited", or "dead"
		case "created":
			status.Created++
		case "running":
			status.Running++
			healthPointer := container.State.Health
			if healthPointer != nil {
				// We perform this check inside the "running" case to avoid counting containers that are not running
				health := *healthPointer
				switch health.Status {
				case "healthy":
					status.RunningHealthy++
				case "unhealthy":
					status.RunningUnhealthy++
				}
			}
		case "paused":
			status.Paused++
		case "restarting":
			status.Restarting++
		case "removing":
			status.Removing++
		case "exited":
			status.Exited++
		case "dead":
			status.Dead++
		}
	}
	containerCount := len(containerStructs)
	if status.Created == containerCount {
		status.Summary = Created
	} else if status.Running == containerCount {
		if status.Running == status.RunningHealthy {
			status.Summary = RunningHealthy
		} else if status.Running == status.RunningUnhealthy {
			status.Summary = RunningUnhealthy
		} else {
			status.Summary = Running
		}
	} else if status.Paused == containerCount {
		status.Summary = Paused
	} else if status.Restarting == containerCount {
		status.Summary = Restarting
	} else if status.Exited == containerCount {
		status.Summary = Exited
	} else if status.Removing == containerCount {
		status.Summary = Removing
	} else if status.Dead == containerCount {
		status.Summary = Dead
	} else {
		status.Summary = Mixed
	}
	return status
}

func extractDependsOn(container types.ContainerJSON) []string {
	var dependsOn []string
	if dependsOnString, ok := container.Config.Labels["com.docker.compose.depends_on"]; ok {
		if dependsOnString != "" {
			dependsOn = strings.Split(dependsOnString, ",")
			for i, service := range dependsOn {
				dependsOn[i] = strings.Split(service, ":")[0]
			}
		}
	}
	return dependsOn
}
