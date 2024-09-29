package main

import (
	"github.com/docker/docker/client"
	"html/template"
	"log"
	"net/http"
)

func HandleRequest(w http.ResponseWriter, r *http.Request, projectName string, tmpl *template.Template, dockerClient *client.Client) {
	switch r.URL.Path {
	case "/":
		ServeFromTemplate(w, projectName, "index", tmpl, dockerClient)
	case "/main":
		ServeFromTemplate(w, projectName, "main", tmpl, dockerClient)
	default:
		http.NotFound(w, r)
	}

}

func ServeFromTemplate(w http.ResponseWriter, projectName string, templateName string, tmpl *template.Template, dockerClient *client.Client) {
	// Get all containers
	projectContainers, err := GetAllContainersInProject(projectName, dockerClient)
	if err != nil {
		reportInternalServerError(w, err, "Unable to list containers", "")
		return
	}

	// Inspect each remaining container
	projectContainerInfos, err := InspectEachContainer(projectContainers, dockerClient)
	if err != nil {
		reportInternalServerError(w, err, "Unable to inspect containers", "")
		return
	}

	// Transform containers to project
	project := TransformContainersToProject(projectContainerInfos, projectName)

	// Render template
	err = tmpl.ExecuteTemplate(w, templateName, project)
	if err != nil {
		reportInternalServerError(w, err, "Unable to render template", "")
	}

	// Log request
	log.Println("Served", templateName, "page.")
}

func reportInternalServerError(w http.ResponseWriter, err error, message string, hint string) {
	log.Println("INTERNAL SERVER ERROR:", message)
	log.Println("Cause:", err.Error())
	if hint != "" {
		log.Println("Hint:", hint)
	}
	http.Error(w, "INTERNAL SERVER ERROR: "+message+"\nCause: "+err.Error(), http.StatusInternalServerError)
	if hint != "" {
		_, _ = w.Write([]byte("Hint: " + hint))
	}
}
