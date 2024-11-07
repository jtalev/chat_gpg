package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type validation_result struct {
	Is_valid	bool
	Msg			string
}

func render_template(w http.ResponseWriter, component string) {
	layout_path := filepath.Join("..", "ui", "layouts", "layout.html")
	nav_path := filepath.Join("..", "ui", "templates", "nav.html")
	tmpl_path := filepath.Join("..", "ui", "views", component+".html")
	tmpl := template.Must(template.ParseFiles(layout_path, nav_path, tmpl_path))
	tmpl.ExecuteTemplate(w, "layout", nil)
}