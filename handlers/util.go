package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type validation_result struct {
	Is_valid bool
	Msg      string
}

func render_template(w http.ResponseWriter, component, title string) {
	// read template file and pass the string
	// version as the component in data struct below
	mainContentPath := filepath.Join("..", "ui", "views", component+".html")
	mainContent, err := os.ReadFile(mainContentPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title     	string
		MainContent template.HTML
	}{
		Title:     title,
		MainContent: template.HTML(mainContent),
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("..", "ui", "layouts", "layout.html"),
		filepath.Join("..", "ui", "templates", "nav.html"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
