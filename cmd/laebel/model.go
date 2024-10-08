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
	Ports       []Port
	DependsOn   []string
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
	Summary          StatusSummary
}

// StatusSummary is a summary of the status of a service. It is used to give a quick overview of a service's status.
type StatusSummary int

const (
	Created StatusSummary = iota
	Running
	// RunningHealthy means the service is running and its health checks are passing.
	RunningHealthy
	// RunningUnhealthy means the service is running, but its health checks are failing.
	RunningUnhealthy
	Paused
	Restarting
	Exited
	Removing
	Dead
	Mixed
)

func (s Status) SummaryIcon() string {
	switch s.Summary {
	case Created:
		return "🆕"
	case Running:
		return "▶️"
	case RunningHealthy:
		return "🟢"
	case RunningUnhealthy:
		return "❌"
	case Paused:
		return "⏸️"
	case Restarting:
		return "🔄"
	case Exited:
		return "⏹️"
	case Removing:
		return "🚮"
	case Dead:
		return "💀"
	case Mixed:
		return "*️⃣"
	default:
		return "❓"
	}
}

func (s Status) SummaryDescription() string {
	switch s.Summary {
	case Created:
		return "The service is created, but not yet running."
	case Running:
		return "The service is running, but has no health information."
	case RunningHealthy:
		return "The service is running and is healthy. 😀"
	case RunningUnhealthy:
		return "The service is running but unhealthy! ☹️"
	case Paused:
		return "The service is paused."
	case Restarting:
		return "The service is restarting."
	case Exited:
		return "The service has exited."
	case Removing:
		return "The service is being removed."
	case Dead:
		return "The service is dead."
	case Mixed:
		return "The service has containers with different states."
	default:
		return "The service has an unknown status."
	}
}

type Link struct {
	Title string // Label: net.henko.laebel.link.<key>.title
	URL   string // Label: net.henko.laebel.link.<key>.url
}

type Port struct {
	Number      string
	Description string // Label: net.henko.laebel.port.<number>.description
}

type Container struct {
	ID      string
	Name    string
	Created string
	Status  string
	Health  string
	Ports   []string
}
