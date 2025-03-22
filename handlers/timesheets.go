package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	application "github.com/jtalev/chat_gpg/application/services"
	"github.com/jtalev/chat_gpg/domain/models"
)

func (h *Handler) ServeTimesheetsView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			component := "timesheets"
			title := "Timesheets - GPG"

			employeeId, err := getEmployeeId(w, r)
			if err != nil {
				log.Println("Unauthorized:", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			timesheetViewData, err := application.InitialTimesheetViewData(employeeId, h.DB)
			if err != nil {
				log.Println("Error retrieving initial view data: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			renderTemplate(w, r, component, title, timesheetViewData)
		},
	)
}

func (h *Handler) GetTimesheetTable() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			employeeId, err := getEmployeeId(w, r)
			if err != nil {
				log.Println("Unauthorized:", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			keys := []string{"arrow", "week_start_date"}
			hxVals, err := parseRequestValues(keys, r)
			if err != nil {
				log.Fatalf("Error parsing request params: ", err)
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			timesheetTableData, err := application.TimesheetTableData(employeeId, hxVals, h.DB)
			if err != nil {
				log.Println("Error retrieving initial view data: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			tmpl, err := template.ParseFiles(
				timesheetTablePath,
				timesheetNavContainerPath,
				timesheetHeadPath,
				existingTimesheetRowPath,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.ExecuteTemplate(w, "timesheetTable", timesheetTableData)
			if err != nil {
				log.Println(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) GetTimesheets() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			outTimesheets, err := application.GetAllTimesheets(h.DB)
			if err != nil {
				log.Println("Error retrieving timesheets from db: ", err)
				http.Error(w, "Internal server error", http.StatusNotFound)
				return
			}
			responseJSON(w, outTimesheets, h.Logger)
		},
	)
}

func (h *Handler) GetTimesheetById() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			requestValues, err := parseRequestValues([]string{"id"}, r)
			if err != nil {
				log.Fatalf("Error parsing query params: ", err)
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			inTimesheetId := requestValues[0]

			outTimesheet, err := application.GetTimesheetById(inTimesheetId, h.DB)
			if err != nil {
				log.Println("Error retrieving timesheet from db: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			responseJSON(w, outTimesheet, h.Logger)
		},
	)
}

func (h *Handler) PutTimesheet() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			keys := []string{"timesheet_id", "time"}
			hxVals, err := parseRequestValues(keys, r)
			if err != nil {
				log.Println("Error parsing request values:", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			timesheetId, time := hxVals[0], hxVals[1]

			outTimesheet, err := application.PutTimesheet(timesheetId, time, h.DB)
			if err != nil {
				log.Println("Error updating timesheet: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			responseJSON(w, outTimesheet, h.Logger)
		},
	)
}

func (h *Handler) InitTimesheetWeek() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("initialising timesheet week")
			employeeId, err := getEmployeeId(w, r)
			if err != nil {
				log.Println(err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			keys := []string{"job_id", "week_start_date"}
			hxVals, err := parseRequestValues(keys, r)
			if err != nil {
				log.Println("Error parsing request vals:", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			jobId := hxVals[0]
			weekStartDate := hxVals[1]

			outTimesheetRow, err := application.InitTimesheetWeek(
				employeeId,
				jobId,
				weekStartDate,
				h.DB)
			if err != nil {
				log.Println("Error initializing timesheet week: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			type Data struct {
				TimesheetRows []application.TimesheetRow
			}
			data := Data{TimesheetRows: outTimesheetRow}

			err = executePartialTemplate(existingTimesheetRowPath, "existingTimesheetRow", data, w)
			if err != nil {
				log.Println("Error executing template:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) GetTimesheetWeekByEmployee() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			employeeId, err := getEmployeeId(w, r)
			if err != nil {
				log.Println(err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			timesheetWeeks, err := application.GetTimesheetWeekByEmployee(employeeId, h.DB)
			if err != nil {
				log.Println("Error retrieving timesheet_week from db: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			responseJSON(w, timesheetWeeks, h.Logger)
		},
	)
}

func (h *Handler) DeleteTimesheetWeek() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			keys := []string{"timesheet_week_id"}
			requestParams, err := parseRequestValues(keys, r)
			if err != nil {
				log.Println("Error parsing request values:", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			id := requestParams[0]

			timesheetWeek, err := application.DeleteTimesheetWeek(id, h.DB)

			log.Println(timesheetWeek)

			fmt.Fprint(w, "<div></div>")
		},
	)
}

func (h *Handler) RenderJobSelectModal() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			type Data struct {
				Jobs          []domain.Job
				WeekStartDate string
			}

			employeeId, err := getEmployeeId(w, r)
			if err != nil {
				log.Println("Unauthorized:", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			hxVals, err := parseRequestValues([]string{"week_start_date"}, r)
			if err != nil {
				log.Println("Error parsing request vals:", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			weekStartDate := hxVals[0]
			jobs, err := application.GetAvailableJobs(employeeId, weekStartDate, h.DB)
			if err != nil {
				log.Println("Error retrieving jobs from database:", err)
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			data := Data{
				Jobs:          jobs,
				WeekStartDate: weekStartDate,
			}

			err = executePartialTemplate(jobSelectModalPath, "jobSelectModal", data, w)
			if err != nil {
				log.Println("Error rendering jobSelectModal:", err)
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				return
			}
		},
	)
}
