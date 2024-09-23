package main

type Project struct {
	Name        string
	Title       string // DOCODASH_PROJECT_TITLE
	Description string // DOCODASH_PROJECT_DESCRIPTION
	Services    []Service
}

type Service struct {
	Name        string
	Title       string // org.opencontainers.image.title
	Description string // org.opencontainers.image.description
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
	Label string
	URL   string
}

type Container struct {
	ID     string
	Names  []string
	Status string
}
