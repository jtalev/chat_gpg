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

var postEmployeeKeys = []string{"id", "employee_id", "first_name", "last_name", "email", "phone_number", "is_admin", "username", "password"}

func (h *Handler) PostEmployee() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			reqVals, err := parseRequestValues(postEmployeeKeys, r)
			if err != nil {
				log.Printf("Error parsing request values: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			employeeDto := application.EmployeeDto{
				ID:          reqVals[0],
				EmployeeId:  reqVals[1],
				FirstName:   reqVals[2],
				LastName:    reqVals[3],
				Email:       reqVals[4],
				PhoneNumber: reqVals[5],
				IsAdmin:     reqVals[6],
				Username:    reqVals[7],
				Password:    reqVals[8],
			}

			employeeDto, err = application.PostEmployee(employeeDto, h.DB)
			if err != nil {
				h.Logger.Errorf("Error posting employee: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			err = executePartialTemplate(adminAddEmployeeModalPath, "adminAddEmployeeModal", employeeDto, w)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

var putEmployeeFormKeys = []string{"id", "employee_id", "first_name", "last_name", "username", "password", "email", "phone"}

func (h *Handler) PutEmployee() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			requestVals, err := parseRequestValues(putEmployeeFormKeys, r)
			if err != nil {
				log.Println("Error parsing request values: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			employeeDto := application.EmployeeDto{}
			employeeDto.ID = requestVals[0]
			employeeDto.EmployeeId = requestVals[1]
			employeeDto.FirstName = requestVals[2]
			employeeDto.LastName = requestVals[3]
			employeeDto.Username = requestVals[4]
			employeeDto.Password = requestVals[5]
			employeeDto.Email = requestVals[6]
			employeeDto.PhoneNumber = requestVals[7]

			employee, err := application.PutEmployee(employeeDto, h.DB)
			if err != nil {
				log.Println("Error updating employee: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			err = executePartialTemplate(adminPutEmployeeModalPath, "adminPutEmployeeModal", employee, w)
			if err != nil {
				log.Println("Error executing template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}
