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
			data, err := application.GetLeaveRequestsForAdmin(h.DB)
			if err != nil {
				log.Printf("Error querying leave database: %v", err)
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			tmpl, err := template.ParseFiles(
				adminLeaveTabPath,
				adminLeaveRequestPath,
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

func (h *Handler) LeaveRequestModal() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			hxvals, err := parseRequestValues([]string{"id"}, r)
			if err != nil {
				log.Printf("Error parsing request values: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			modalData, err := application.GetLeaveRequestByIdForAdmin(hxvals[0], h.DB)
			if err != nil {
				log.Printf("Error getting admin leave modal data: %v", err)
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			err = executePartialTemplate(adminLeaveRequestModalPath, "adminLeaveModal", modalData, w)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) LeaveRequestFinalise() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			hxvals, err := parseRequestValues([]string{"id", "approved"}, r)
			if err != nil {
				log.Printf("Error parsing request values: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			id, isApproved := hxvals[0], hxvals[1]
			_, err = application.AdminUpdateLeaveRequest(id, isApproved, h.DB)
			if err != nil {
				log.Printf("Error finalising leave request: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			data, err := application.GetLeaveRequestsForAdmin(h.DB)
			if err != nil {
				log.Printf("Error querying leave database: %v", err)
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			tmpl, err := template.ParseFiles(
				adminLeaveTabPath,
				adminLeaveRequestPath,
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

func (h *Handler) AddEmployeeModal() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := executePartialTemplate(adminAddEmployeeModalPath, "adminAddEmployeeModal", nil, w)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}
