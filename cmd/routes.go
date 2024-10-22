package main

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
)

func add_routes(mux *http.ServeMux, ctx context.Context) {
	fileServer := http.FileServer(http.Dir("../static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/", serve_index())
}

func serve_index() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			path := filepath.Join("..", "static", "index.html")
			tmpl := template.Must(template.ParseFiles(path))
			if err := tmpl.Execute(w, nil); err != nil {
				fmt.Printf("error executing template: %v", err)
				http.Error(w, "error executing tempate", http.StatusInternalServerError)
			}
		},
	)
}
