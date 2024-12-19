package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/jtalev/chat_gpg/repository"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

func responseJson(w http.ResponseWriter, data any, sugar *zap.SugaredLogger) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		sugar.Errorf("Error encoding leave requests: %v", err)
		http.Error(w, "failed to fetch leave requests", http.StatusInternalServerError)
	}
}

func renderTemplate(
	w http.ResponseWriter,
	r *http.Request,
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

	isAdmin, ok := r.Context().Value("is_admin").(bool)
	if !ok {
		http.Error(w, "unable to retrieve is_admin", http.StatusUnauthorized)
		return
	}

	data := struct {
		Title     string
		Component string
		IsAdmin   bool
		Data      interface{}
	}{
		Title:     title,
		Component: component,
		IsAdmin:   isAdmin,
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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

type Handler struct {
	DB *sql.DB
	Store *sessions.CookieStore
	Sugar *zap.SugaredLogger
}

func (h *Handler) ServeLoginView(w http.ResponseWriter, r *http.Request) {
	login_path := filepath.Join("..", "ui", "views", "login.html")
	tmpl := template.Must(template.ParseFiles(login_path))
	tmpl.Execute(w, nil)
}

func (h *Handler) ServeAccountView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data := getAccountData()
			component := "account"
			title := "Account - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}

func (h *Handler) ServeAdminView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data := getAdminData()
			component := "admin"
			title := "Admin - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}

func (h *Handler) ServeDashboardView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data := getDashboardData()
			component := "dashboard"
			title := "Dashboard - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}

func (h *Handler) ServeJobsView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data := getJobsData()
			component := "jobs"
			title := "Jobs - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}

func (h *Handler) ServeLeaveView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data, err := repository.GetLeaveRequestsByEmployee(w, r, h.DB, h.Sugar)
			if err != nil {
				h.Sugar.Errorf("Error getting leave page data: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			fmt.Println(data)
			component := "leave"
			title := "Leave - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}

func (h *Handler) ServeTimesheetsView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data := getTimesheetData()
			component := "timesheets"
			title := "Timesheets - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}
