package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"reflect"
)

func RenderDocument(w http.ResponseWriter, err error, project Project) {
	// Load template
	tmpl, err := template.ParseFiles(
		filepath.Join("web", "templates", "index.html"),
		filepath.Join("web", "templates", "serviceStatus.html"),
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

// isArrayOrSlice checks if the passed value is an array or slice.
func isArrayOrSlice(v interface{}) bool {
	rv := reflect.ValueOf(v)
	return rv.Kind() == reflect.Array || rv.Kind() == reflect.Slice
}

// isMap checks if the passed value is a map.
func isMap(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.Map
}

// isStruct checks if the passed value is a struct.
func isStruct(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.Struct
}

// structFields returns a map of field names and values for a struct.
func structFields(v interface{}) map[string]interface{} {
	rv := reflect.ValueOf(v)
	rt := rv.Type()
	fields := make(map[string]interface{})

	// Loop over struct fields and add them to the map
	for i := 0; i < rv.NumField(); i++ {
		fieldName := rt.Field(i).Name
		fieldValue := rv.Field(i).Interface()
		fields[fieldName] = fieldValue
	}
	return fields
}
