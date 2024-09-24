package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func RenderDocument(w http.ResponseWriter, err error, project Project) {
	// Load template
	tmpl, err := template.ParseFiles(
		filepath.Join("web", "templates", "index.html"),
		filepath.Join("web", "templates", "serviceStatus.html"),
		filepath.Join("web", "templates", "clipboard.html"),
	)
	if err != nil {
		InternalServerError(w, err, "Unable to load template", "")
		return
	}

	// Render template
	err = tmpl.Execute(w, project)
	if err != nil {
		InternalServerError(w, err, "Unable to render template", "")
	}
}
