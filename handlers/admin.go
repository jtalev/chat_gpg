package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

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

func (h *Handler) RenderEmployeeTab() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			employees, err := infrastructure.GetEmployees(h.DB)
			if err != nil {
				log.Printf("Error querying employee database: %v", err)
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}

			data := struct {
				Employees []domain.Employee
			}{
				Employees: employees,
			}

			tmpl, err := template.ParseFiles(
				adminEmployeeTabPath,
				adminEmployeeListPath,
			)
			if err != nil {
				log.Printf("Error parsing templates: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.ExecuteTemplate(w, "adminEmployeeTab", data)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
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
				log.Printf("Error parsing templates: %v", err)
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
			jobDto := application.JobDto{}
			jobDto.ID = "-1"
			err := executePartialTemplate(addJobModalPath, "addJobModal", jobDto, w)
			if err != nil {
				log.Printf("Error executing addJobModal.html: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

var putJobModalKeys = []string{"id"}

func (h *Handler) PutJobModal() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			rVals, err := parseRequestValues(putJobModalKeys, r)
			if err != nil {
				log.Printf("Error parsing request values: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(rVals[0])
			if err != nil {
				log.Printf("Error converting string id to int: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			job, err := application.GetJobById(id, h.DB)
			if err != nil {
				log.Printf("Error getting job: %v", err)
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			isAvailable := "true"
			if !job.IsComplete {
				isAvailable = "false"
			}

			jobDto := application.JobDto{
				ID:         strconv.Itoa(job.ID),
				Name:       job.Name,
				Number:     strconv.Itoa(job.Number),
				Address:    job.Address,
				Suburb:     job.Suburb,
				PostCode:   job.PostCode,
				City:       job.City,
				IsComplete: isAvailable,
			}

			err = executePartialTemplate(putJobModalPath, "putJobModal", jobDto, w)
			if err != nil {
				log.Printf("Error executing putJobModal.html: %v", err)
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

func (h *Handler) PutEmployeeModal() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			keys := []string{"id", "employee_id"}
			vals, err := parseRequestValues(keys, r)
			if err != nil {
				log.Println("Error parsing request values: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(vals[0])
			if err != nil {
				log.Println("Error converting request value to int: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			employee, err := application.GetEmployeeById(id, h.DB)
			if err != nil {
				log.Println("Error getting employee: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			employeeAuth, err := infrastructure.GetEmployeeAuthByEmployeeId(vals[1], h.DB)
			if err != nil {
				log.Println("Error getting employee auth: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			isAdmin := "false"
			if employee.IsAdmin == true {
				isAdmin = "true"
			}

			employeeDto := application.EmployeeDto{
				ID:          vals[0],
				EmployeeId:  vals[1],
				FirstName:   employee.FirstName,
				LastName:    employee.LastName,
				Username:    employeeAuth.Username,
				Email:       employee.Email,
				PhoneNumber: employee.PhoneNumber,
				IsAdmin:     isAdmin,
			}

			err = executePartialTemplate(adminPutEmployeeModalPath, "adminPutEmployeeModal", employeeDto, w)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) RenderSafetyTab() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := executePartialTemplate(adminSafetyTabPath, "adminSafetyTab", nil, w)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) AdminServeIncidentReportContent() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			incidentReports, err := infrastructure.GetIncidentReports(h.DB)
			if err != nil {
				log.Printf("Error getting incident reports: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}

			err = executePartialTemplate(adminIncidentReportListPath, "adminIncidentReportList", incidentReports, w)
			if err != nil {
				log.Printf("error executing adminIncidentReportList.html: %v", err)
				http.Error(w, "error executing adminIncidentReportList.html, internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) AdminServeSwmListContent() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := executePartialTemplate(adminSwmListPath, "adminSwmList", []int{1, 2, 3, 4, 5}, w)
			if err != nil {
				log.Printf("error executing adminSwmList.html: %v", err)
				http.Error(w, "error executing adminSwmList.html, internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}
