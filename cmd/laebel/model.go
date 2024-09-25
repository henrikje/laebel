package main

type Project struct {
	Name          string // Label: com.docker.compose.project
	Title         string // Env: LAEBEL_PROJECT_TITLE
	Description   string // Env: LAEBEL_PROJECT_DESCRIPTION
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
	Created    int
	Running    int
	Paused     int
	Restarting int
	Exited     int
	Removing   int
	Dead       int
	Stopped    int
}

type Link struct {
	Label string // Label: net.henko.laebel.link.<key>.label
	URL   string // Label: net.henko.laebel.link.<key>.url
}

type Container struct {
	ID     string
	Names  []string
	Status string
}
