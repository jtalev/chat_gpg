package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type validation_result struct {
	Is_valid bool
	Msg      string
}

func render_template(w http.ResponseWriter, component, title string) {
	data := struct {
		Title			string
		Component	string
	}{
		Title: 			title,
		Component: 	component,
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("..", "ui", "layouts", "layout.html"),
		filepath.Join("..", "ui", "templates", "nav.html"),
		filepath.Join("..", "ui", "views", component+".html"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
