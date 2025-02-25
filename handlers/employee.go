package handlers

import (
	"log"
	"net/http"

	application "github.com/jtalev/chat_gpg/application/services"
	domain "github.com/jtalev/chat_gpg/domain/models"
)

type EmployeeData struct {
	Employees []domain.Employee
}

func (h *Handler) DeleteEmployee() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			isAdmin, err := getIsAdmin(w, r)
			if err != nil {
				log.Printf("Error getting admin cookie: %v", err)
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			if !isAdmin {
				log.Printf("Unauthorized: isAdmin = %v", isAdmin)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			keys := []string{"id"}
			hxVals, err := parseRequestValues(keys, r)
			if err != nil {
				log.Printf("Error parsing request values: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			employees, err := application.DeleteAndReturnEmployees(hxVals[0], h.DB)
			data := EmployeeData{employees}
			err = executePartialTemplate(adminEmployeeListPath, "adminEmployeeList", data, w)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) PostEmployee() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			isAdmin, err := getIsAdmin(w, r)
			if err != nil {
				log.Printf("Error getting admin cookie: %v", err)
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			if !isAdmin {
				log.Printf("Unauthorized: isAdmin = %v", isAdmin)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			keys := []string{"employee_id", "first_name", "last_name", "email", "phone_number", "is_admin", "username", "password"}
			hxVals, err := parseRequestValues(keys, r)
			if err != nil {
				log.Printf("Error parsing request values: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			employees, err := application.PostAndReturnEmployees(hxVals, h.DB)
			data := EmployeeData{employees}
			err = executePartialTemplate(adminEmployeeListPath, "adminEmployeeList", data, w)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}
