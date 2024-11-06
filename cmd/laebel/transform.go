package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/go-connections/nat"
	"maps"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

func TransformContainersToProject(projectContainers []types.ContainerJSON, projectVolumes []*volume.Volume, projectNetworks []network.Summary, projectName string) Project {
	containers := filterOutHiddenContainers(projectContainers)
	containersByServiceName := groupContainersByServiceName(containers)
	volumesByName := extractVolumes(projectVolumes, containersByServiceName)
	networksByID := extractNetworks(projectNetworks, containersByServiceName)
	services := createServiceForEachServiceName(containersByServiceName, volumesByName, networksByID)
	serviceGroups := groupServicesByGroup(containers, services)
	projectLinks := extractProjectLinks()
	volumes := slices.Collect(maps.Values(volumesByName))
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Name < volumes[j].Name
	})
	networks := slices.Collect(maps.Values(networksByID))
	sort.Slice(networks, func(i, j int) bool {
		return networks[i].Name < networks[j].Name
	})
	return Project{
		Name:          projectName,
		Title:         os.Getenv("LAEBEL_PROJECT_TITLE"),
		Description:   os.Getenv("LAEBEL_PROJECT_DESCRIPTION"),
		Links:         projectLinks,
		Logo:          os.Getenv("LAEBEL_PROJECT_LOGO"),
		ServiceGroups: serviceGroups,
		Volumes:       volumes,
		Networks:      networks,
	}
}

func extractVolumes(projectVolumes []*volume.Volume, containersByServiceName map[string][]types.ContainerJSON) map[string]Volume {
	volumes := make(map[string]Volume)
	for _, projectVolume := range projectVolumes {
		if projectVolume.Labels["net.henko.laebel.hidden"] == "true" {
			continue
		}
		volumes[projectVolume.Name] = Volume{
			Name:        projectVolume.Labels["com.docker.compose.volume"],
			Title:       projectVolume.Labels["net.henko.laebel.title"],
			Description: projectVolume.Labels["net.henko.laebel.description"],
			Driver:      projectVolume.Driver,
			Services:    servicesUsingVolume(projectVolume.Name, containersByServiceName),
		}
	}
	return volumes
}

func servicesUsingVolume(volumeName string, containersByServiceName map[string][]types.ContainerJSON) []string {
	services := make([]string, 0)
	for serviceName, serviceContainers := range containersByServiceName {
		for _, container := range serviceContainers {
			for _, containerMount := range container.Mounts {
				if containerMount.Type == "volume" && containerMount.Name == volumeName {
					services = append(services, serviceName)
					break
				}
			}
		}
	}
	return slices.Compact(slices.Sorted(slices.Values(services)))
}

func extractNetworks(projectNetworks []network.Summary, containersByServiceName map[string][]types.ContainerJSON) map[string]Network {
	networks := make(map[string]Network)
	for _, projectNetwork := range projectNetworks {
		if projectNetwork.Labels["net.henko.laebel.hidden"] == "true" {
			continue
		}
		networks[projectNetwork.ID] = Network{
			Name:        projectNetwork.Labels["com.docker.compose.network"],
			Title:       projectNetwork.Labels["net.henko.laebel.title"],
			Description: projectNetwork.Labels["net.henko.laebel.description"],
			Driver:      projectNetwork.Driver,
			Services:    servicesUsingNetwork(projectNetwork.Name, containersByServiceName),
		}
	}
	if len(networks) == 1 && networks[slices.Collect(maps.Keys(networks))[0]].Name == "default" {
		// Remove default network if it's the only one
		return nil
	}
	return networks
}

func servicesUsingNetwork(networkName string, containersByServiceName map[string][]types.ContainerJSON) []string {
	services := make([]string, 0)
	for serviceName, serviceContainers := range containersByServiceName {
		for _, container := range serviceContainers {
			for containerNetworkName := range container.NetworkSettings.Networks {
				if containerNetworkName == networkName {
					services = append(services, serviceName)
					break
				}
			}
		}
	}
	return slices.Compact(slices.Sorted(slices.Values(services)))
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

func createServiceForEachServiceName(containersByServiceName map[string][]types.ContainerJSON, volumesByName map[string]Volume, networksByID map[string]Network) []Service {
	var services []Service
	for serviceName, serviceContainers := range containersByServiceName {
		service := transformContainersToService(serviceContainers, serviceName, volumesByName, networksByID)
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

func transformContainersToService(serviceContainers []types.ContainerJSON, serviceName string, name map[string]Volume, byName map[string]Network) Service {
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
	var volumes []Volume
	for _, containerMount := range container.Mounts {
		if containerMount.Type == "volume" {
			projectVolume, found := name[containerMount.Name]
			if found {
				volumes = append(volumes, projectVolume)
			}
		}
	}
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Name < volumes[j].Name
	})
	var networks []Network
	for _, containerNetwork := range container.NetworkSettings.Networks {
		projectNetwork, found := byName[containerNetwork.NetworkID]
		if found {
			networks = append(networks, projectNetwork)
		}
	}
	sort.Slice(networks, func(i, j int) bool {
		return networks[i].Name < networks[j].Name
	})
	service := Service{
		Name:        serviceName,
		Title:       container.Config.Labels["org.opencontainers.image.title"],
		Description: container.Config.Labels["org.opencontainers.image.description"],
		Image:       container.Config.Image,
		Status:      status,
		Links:       links,
		Ports:       extractServicePorts(serviceContainers),
		Volumes:     volumes,
		Networks:    networks,
		DependsOn:   dependsOn,
		Containers:  containerStructs,
	}
	return service
}

func extractServicePorts(containers []types.ContainerJSON) []Port {
	ports := make([]Port, 0)
	for _, container := range containers {
		for exposedPort, portBindings := range container.HostConfig.PortBindings {
			for _, portBinding := range portBindings {
				port := Port{
					Number:      portBinding.HostPort,
					Description: container.Config.Labels["net.henko.laebel.port."+exposedPort.Port()+".description"],
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
