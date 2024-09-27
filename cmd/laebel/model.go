package main

type Project struct {
	Name          string // Label: com.docker.compose.project
	Title         string // Env: LAEBEL_PROJECT_TITLE
	Description   string // Env: LAEBEL_PROJECT_DESCRIPTION
	Links         []Link // Env: LAEBEL_PROJECT_URL, LAEBEL_PROJECT_DOCUMENTATION, LAEBEL_PROJECT_SOURCE
	Icon          string // Env: LAEBEL_PROJECT_ICON
	ServiceGroups []ServiceGroup
}

type ServiceGroup struct {
	Name     string // Label: net.henko.laebel.group
	Services []Service
}

type Service struct {
	Name        string // Label: com.docker.compose.service
	Title       string // Label: org.opencontainers.image.title
	Description string // Label: org.opencontainers.image.description
	Image       string
	Status      Status
	Links       []Link
	DependsOn   []string // Label: com.docker.compose.depends_on
	Containers  []Container
}

type Status struct {
	Created          int
	Running          int
	RunningHealthy   int
	RunningUnhealthy int
	Paused           int
	Restarting       int
	Exited           int
	Removing         int
	Dead             int
	Stopped          int
	Summary          StatusSummary
}

// StatusSummary is a summary of the status of a service. It is used to give a quick overview of a service's status.
type StatusSummary int

const (
	// Running means the service is running, but we do not have any information about its health.
	Running StatusSummary = iota
	// RunningHealthy means the service is running and its health checks are passing.
	RunningHealthy
	// RunningUnhealthy means the service is running, but its health checks are failing.
	RunningUnhealthy
	// NotRunning means the service is not running, because it is paused, stopped, have exited, or is dead.
	NotRunning
	// Other means the service is in a state not covered by the other statuses.
	Other
)

type Link struct {
	Label string // Label: net.henko.laebel.link.<key>.label
	URL   string // Label: net.henko.laebel.link.<key>.url
}

type Container struct {
	ID     string
	Name   string
	Status string
}
