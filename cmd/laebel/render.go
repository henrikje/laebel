package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

var escapeReplacer = strings.NewReplacer(".", "#46;", "(", "#40;", ")", "#41;")

func RenderDocument(w http.ResponseWriter, err error, project Project) {
	// Load template
	tmpl, err := template.New("index.html").Funcs(
		template.FuncMap{
			"escape": func(raw string) string { return escapeReplacer.Replace(raw) },
		}).ParseFiles(
		filepath.Join("web", "templates", "index.html"),
		filepath.Join("web", "templates", "serviceGraph.html"),
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
