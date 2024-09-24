package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

func RenderDocument(w http.ResponseWriter, err error, project Project) {
	// Load template
	tmpl, err := template.New("index.html").Funcs(
		template.FuncMap{
			"escape": func(raw string) string {
				// TODO Use a markdown escape function instead
				return strings.Replace(raw, ".", "#46;", -1)
			},
		}).ParseFiles(
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
