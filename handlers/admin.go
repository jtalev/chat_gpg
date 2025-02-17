package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	application "github.com/jtalev/chat_gpg/application/services"
	domain "github.com/jtalev/chat_gpg/domain/models"
	infrastructure "github.com/jtalev/chat_gpg/infrastructure/repository"
)

type AdminData struct {
	Employees []domain.Employee
}

func getInitialAdminData(db *sql.DB) (AdminData, error) {
	employees, err := application.GetEmployees(db)
	if err != nil {
		return AdminData{}, err
	}
	data := AdminData{
		employees,
	}
	return data, nil
}

func (h *Handler) ServeAdminView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data, err := getInitialAdminData(h.DB)
			if err != nil {
				log.Printf("Error getting initial admin data: %v", err)
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			component := "admin"
			title := "Admin - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}

func (h *Handler) RenderJobTab() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data, err := infrastructure.GetJobs(h.DB)
			if err != nil {
				log.Printf("Error querying job database: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			tmpl, err := template.ParseFiles(
				adminJobTabPath,
				adminJobListPath,
			)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.ExecuteTemplate(w, "adminJobTab", data)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) RenderJobList() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data, err := infrastructure.GetJobs(h.DB)
			if err != nil {
				log.Printf("Error querying job database: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			err = executePartialTemplate(adminJobListPath, "adminJobList", data, w)
			if err != nil {
				log.Printf("Error executing adminJobList.html: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) AddJobModal() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			jobs, err := application.GetJobs(h.DB)
			if err != nil {
				log.Printf("Error getting job data: %v", err)
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			err = executePartialTemplate(addJobModalPath, "addJobModal", jobs, w)
			if err != nil {
				log.Printf("Error executing addJobModal.html: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) RenderLeaveTab() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data, err := infrastructure.GetLeaveRequests(h.DB)
			if err != nil {
				log.Printf("Error querying leave database: %v", err)
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			tmpl, err := template.ParseFiles(
				adminLeaveTabPath,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.ExecuteTemplate(w, "adminLeaveTab", data)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}
