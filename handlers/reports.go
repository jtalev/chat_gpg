package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	report "github.com/jtalev/chat_gpg/application/services/report"
)

type ReportsViewData struct {
	TimesheetReportData report.TimesheetReportData
}

func (h *Handler) ServeReportsView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			initialTimesheetReportData, err := report.InitialTimesheetReportData(h.DB)
			if err != nil {
				log.Println("Error getting initial timesheet report data:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			outData := ReportsViewData{
				TimesheetReportData: initialTimesheetReportData,
			}
			component := "reports"
			title := "Reports - GPG"
			renderTemplate(w, r, component, title, outData)
		},
	)
}

func (h *Handler) InitTimesheetReport() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			initialTimesheetReportData, err := report.InitialTimesheetReportData(h.DB)
			if err != nil {
				log.Printf("error initializing timesheet report data: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			err = executePartialTemplate(timesheetReportPath, "timesheetReport", initialTimesheetReportData, w)
			if err != nil {
				log.Printf("error executing partial template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) GetEmployeeTimesheetReport() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			isAdmin, err := getIsAdmin(w, r)
			if err != nil || !isAdmin {
				log.Println("Unauthorized user:", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			keys := []string{"id", "week_start_date"}
			requestVals, err := parseRequestValues(keys, r)
			if err != nil {
				log.Println("Error parsing request values:", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			employeeId, weekStartDate := requestVals[0], requestVals[1]

			outData, err := report.GetEmployeeTimesheetReport(employeeId, weekStartDate, h.DB)
			if err != nil {
				log.Println("Error getting employee timesheet report:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			tmpl, err := template.ParseFiles(employeeTimesheetReportPath, reportEmployeeLeaveRequestPath)
			if err != nil {
				log.Println("Error parsing file:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			err = tmpl.ExecuteTemplate(w, "employeeTimesheetReport", outData)
			if err != nil {
				log.Println("Error executing template:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) PrevEmployeeTimesheetReport() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			isAdmin, err := getIsAdmin(w, r)
			if err != nil || !isAdmin {
				log.Println("Unauthorized user:", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			keys := []string{"id", "week_start_date"}
			requestVals, err := parseRequestValues(keys, r)
			if err != nil {
				log.Println("Error parsing request values:", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			employeeId, weekStartDate := requestVals[0], requestVals[1]

			outData, err := report.GetPrevEmployeeTimesheetReport(employeeId, weekStartDate, h.DB)
			if err != nil {
				log.Println("Error getting employee timesheet report:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			tmpl, err := template.ParseFiles(employeeTimesheetReportPath, reportEmployeeLeaveRequestPath)
			if err != nil {
				log.Println("Error parsing file:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			err = tmpl.ExecuteTemplate(w, "employeeTimesheetReport", outData)
			if err != nil {
				log.Println("Error executing template:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) NextEmployeeTimesheetReport() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			isAdmin, err := getIsAdmin(w, r)
			if err != nil || !isAdmin {
				log.Println("Unauthorized user:", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			keys := []string{"id", "week_start_date"}
			requestVals, err := parseRequestValues(keys, r)
			if err != nil {
				log.Println("Error parsing request values:", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			employeeId, weekStartDate := requestVals[0], requestVals[1]

			outData, err := report.GetNextEmployeeTimesheetReport(employeeId, weekStartDate, h.DB)
			if err != nil {
				log.Println("Error getting employee timesheet report:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			tmpl, err := template.ParseFiles(employeeTimesheetReportPath, reportEmployeeLeaveRequestPath)
			if err != nil {
				log.Println("Error parsing file:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			err = tmpl.ExecuteTemplate(w, "employeeTimesheetReport", outData)
			if err != nil {
				log.Println("Error executing template:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) InitJobReportData() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			initData, err := report.InitJobReportData(h.DB)
			if err != nil {
				log.Printf("error getting initial job report data: %v", err)
				http.Error(w, "internal server error: error getting initial job report data", http.StatusInternalServerError)
				return
			}
			err = executePartialTemplate(initJobReportPath, "initJobReport", initData, w)
			if err != nil {
				log.Printf("error executing partial template: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

type JobSelect struct {
	ID string `json:"id"`
}

func (h *Handler) GetJobReport() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			jobSelect := JobSelect{}
			if err := json.NewDecoder(r.Body).Decode(&jobSelect); err != nil {
				log.Printf("error decoding json: %v", err)
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			id, err := strconv.Atoi(jobSelect.ID)
			if err != nil {
				log.Printf("error converting id to int, bad request: %v", err)
				http.Error(w, "error converting job id to int: bad request", http.StatusBadRequest)
				return
			}

			out, err := report.GetJobReportData(id, 1, h.DB)
			if err != nil {
				log.Printf("error getting job report data: %v", err)
				http.Error(w, "internal server error: error getting job report data from server", http.StatusInternalServerError)
				return
			}

			err = executePartialTemplate(jobReportPath, "jobReport", out, w)
			if err != nil {
				log.Printf("error executing job report templates: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}
