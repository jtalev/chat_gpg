package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/jtalev/chat_gpg/auth"
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
	store *sessions.CookieStore,
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

	session, err := store.Get(r, "employee_session")
	if err != nil {
		fmt.Errorf("Error getting store: %v", err)
		http.Error(w, "Error getting session", http.StatusInternalServerError)
		return
	}
	isAdminValue := session.Values["is_admin"]
	isAdmin := false
	if isAdminValue == true {
		isAdmin = true
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ServeAccountView(store *sessions.CookieStore, sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			auth.RedirectUnauthorisedUser(w, r, store, sugar)
			data := getAccountData()
			component := "account"
			title := "Account - GPG"
			renderTemplate(w, r, store, component, title, data)
		},
	)
}

func ServeAdminView(store *sessions.CookieStore, sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			auth.RedirectUnauthorisedUser(w, r, store, sugar)
			data := getAdminData()
			component := "admin"
			title := "Admin - GPG"
			renderTemplate(w, r, store, component, title, data)
		},
	)
}

func ServeDashboardView(store *sessions.CookieStore, sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			auth.RedirectUnauthorisedUser(w, r, store, sugar)
			data := getDashboardData()
			component := "dashboard"
			title := "Dashboard - GPG"
			renderTemplate(w, r, store, component, title, data)
		},
	)
}

func ServeJobsView(store *sessions.CookieStore, sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			auth.RedirectUnauthorisedUser(w, r, store, sugar)
			data := getJobsData()
			component := "jobs"
			title := "Jobs - GPG"
			renderTemplate(w, r, store, component, title, data)
		},
	)
}

func ServeLeaveView(store *sessions.CookieStore, sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			auth.RedirectUnauthorisedUser(w, r, store, sugar)
			data := GetLeaveRequestById(sugar)
			component := "leave"
			title := "Leave - GPG"
			renderTemplate(w, r, store, component, title, data)
		},
	)
}

func ServeLoginView(w http.ResponseWriter, r *http.Request) {
	login_path := filepath.Join("..", "ui", "views", "login.html")
	tmpl := template.Must(template.ParseFiles(login_path))
	tmpl.Execute(w, nil)
}

func ServeTimesheetsView(store *sessions.CookieStore, sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			auth.RedirectUnauthorisedUser(w, r, store, sugar)
			data := getTimesheetData()
			component := "timesheets"
			title := "Timesheets - GPG"
			renderTemplate(w, r, store, component, title, data)
		},
	)
}
