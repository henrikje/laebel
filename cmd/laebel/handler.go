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
		ServeFromProjectTemplate(w, projectName, "index", tmpl, dockerClient)
	case "/service":
		serviceName := r.URL.Query().Get("name")
		if serviceName == "" {
			reportBadRequest(w, "No service specified", "Specify a service name in the 'name' query parameter.")
			return
		}
		ServeFromServiceTemplate(w, projectName, serviceName, tmpl, dockerClient)
	case "/graph-status.css":
		ServeFromProjectTemplate(w, projectName, "serviceGraphStatus", tmpl, dockerClient)
	default:
		http.NotFound(w, r)
	}

}

func ServeFromProjectTemplate(w http.ResponseWriter, projectName string, templateName string, tmpl *template.Template, dockerClient *client.Client) {
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

	// Get all volumes
	projectVolumes, err := GetAllVolumesInProject(projectName, dockerClient)
	if err != nil {
		reportInternalServerError(w, err, "Unable to list volumes", "")
		return
	}

	// Get all networks
	projectNetworks, err := GetAllNetworksInProject(projectName, dockerClient)
	if err != nil {
		reportInternalServerError(w, err, "Unable to list networks", "")
		return
	}

	// Transform containers to project
	project := TransformContainersToProject(projectContainerInfos, projectVolumes, projectNetworks, projectName)

	// Render template
	err = tmpl.ExecuteTemplate(w, templateName, project)
	if err != nil {
		reportInternalServerError(w, err, "Unable to render template", "")
	}
}

func ServeFromServiceTemplate(w http.ResponseWriter, projectName string, serviceName string, tmpl *template.Template, dockerClient *client.Client) {
	// TODO Avoid reading the whole project just to get one service

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

	// Get all volumes
	projectVolumes, err := GetAllVolumesInProject(projectName, dockerClient)
	if err != nil {
		reportInternalServerError(w, err, "Unable to list volumes", "")
		return
	}

	// Get all networks
	projectNetworks, err := GetAllNetworksInProject(projectName, dockerClient)
	if err != nil {
		reportInternalServerError(w, err, "Unable to list networks", "")
		return
	}

	// Transform containers to project
	project := TransformContainersToProject(projectContainerInfos, projectVolumes, projectNetworks, projectName)

	// Select service
	var activeService Service
	for _, serviceGroup := range project.ServiceGroups {
		for _, service := range serviceGroup.Services {
			if service.Name == serviceName {
				activeService = service
			}
		}
	}

	// Render template
	err = tmpl.ExecuteTemplate(w, "service", activeService)
	if err != nil {
		reportInternalServerError(w, err, "Unable to render template", "")
	}
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

func reportBadRequest(w http.ResponseWriter, message string, hint string) {
	log.Println("BAD REQUEST:", message)
	if hint != "" {
		log.Println("Hint:", hint)
	}
	http.Error(w, "BAD REQUEST: "+message, http.StatusBadRequest)
	if hint != "" {
		_, _ = w.Write([]byte("Hint: " + hint))
	}
}
