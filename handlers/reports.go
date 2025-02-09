package handlers

import (
	"log"
	"net/http"

	"github.com/jtalev/chat_gpg/application/services"
)

type ReportsViewData struct {
	TimesheetReportData application.TimesheetReportData
}

func (h *Handler) ServeReportsView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			initialTimesheetReportData, err := application.InitialTimesheetReportData(h.DB)
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

			outData, err := application.GetEmployeeTimesheetReport(employeeId, weekStartDate, h.DB)
			if err != nil {
				log.Println("Error getting employee timesheet report:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			err = executePartialTemplate(employeeTimesheetReportPath, "employeeTimesheetReport", outData, w)
			if err != nil {
				log.Println("Error rendering template:", err)
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

			outData, err := application.GetPrevEmployeeTimesheetReport(employeeId, weekStartDate, h.DB)
			if err != nil {
				log.Println("Error getting employee timesheet report:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			err = executePartialTemplate(employeeTimesheetReportPath, "employeeTimesheetReport", outData, w)
			if err != nil {
				log.Println("Error rendering template:", err)
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

			outData, err := application.GetNextEmployeeTimesheetReport(employeeId, weekStartDate, h.DB)
			if err != nil {
				log.Println("Error getting employee timesheet report:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			err = executePartialTemplate(employeeTimesheetReportPath, "employeeTimesheetReport", outData, w)
			if err != nil {
				log.Println("Error rendering template:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}
