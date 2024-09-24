package main

type Project struct {
	Name          string // Label: com.docker.compose.project
	Title         string // Env: DOCODASH_PROJECT_TITLE
	Description   string // Env: DOCODASH_PROJECT_DESCRIPTION
	ServiceGroups []ServiceGroup
}

type ServiceGroup struct {
	Name     string // Label: net.henko.docodash.group
	Services []Service
}

type Service struct {
	Name        string // Label: com.docker.compose.service
	Title       string // Label: org.opencontainers.image.title
	Description string // Label: org.opencontainers.image.description
	Image       string
	Status      Status
	Links       []Link
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
	Label string // Label: net.henko.docodash.<key>.label
	URL   string // Label: net.henko.docodash.<key>.url
}

type Container struct {
	ID     string
	Names  []string
	Status string
}
