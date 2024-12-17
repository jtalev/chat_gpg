package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"go.uber.org/zap"
)

func renderTemplate(
	w http.ResponseWriter,
	component, title string,
	componentData interface{},
) {
	layoutPath := filepath.Join("..", "ui", "layouts", "layout.html")
	navPath := filepath.Join("..", "ui", "templates", "nav.html")
	dashboardPath := filepath.Join("..", "ui", "views", "dashboard.html")
	jobsPath := filepath.Join("..", "ui", "views", "jobs.html")
	timesheetsPath := filepath.Join("..", "ui", "views", "timesheets.html")
	leavePath := filepath.Join("..", "ui", "views", "leave.html")
	adminPath := filepath.Join("..", "ui", "views", "admin.html")
	accountPath := filepath.Join("..", "ui", "views", "account.html")

	data := struct {
		Title     string
		Component string
		Data      interface{}
	}{
		Title:     title,
		Component: component,
		Data:      componentData,
	}

	tmpl, err := template.ParseFiles(
		layoutPath,
		navPath,
		dashboardPath,
		jobsPath,
		timesheetsPath,
		leavePath,
		adminPath,
		accountPath,
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

func responseJson(w http.ResponseWriter, data any, sugar *zap.SugaredLogger) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		sugar.Errorf("Error encoding leave requests: %v", err)
		http.Error(w, "failed to fetch leave requests", http.StatusInternalServerError)
	}
}
