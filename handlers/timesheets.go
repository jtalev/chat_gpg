package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"github.com/jtalev/chat_gpg/services"
)

func parseRequestValues(keys []string, r *http.Request) ([]string, error) {
	out := make([]string, len(keys))

	err := r.ParseForm()
	if err != nil {
		return out, err
	}

	for i := range keys {
		val := r.FormValue(keys[i])
		out[i] = val
	}

	return out, nil
}

func decodeJSON(payload interface{}, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(payload)
	return err
}

func (h *Handler) ServeTimesheetsView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			component := "timesheets"
			title := "Timesheets - GPG"

			timesheetViewData, err := services.InitialTimesheetViewData(h.DB)
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

			keys := []string{"arrow", "week_start_date"}
			hxVals, err := parseRequestValues(keys, r)
			if err != nil {
				log.Fatalf("Error parsing request params: ", err)
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			timesheetTableData, err := services.TimesheetTableData(hxVals, h.DB)
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
			outTimesheets, err := services.GetAllTimesheets(h.DB)
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

			outTimesheet, err := services.GetTimesheetById(inTimesheetId, h.DB)
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

			outTimesheet, err := services.PutTimesheet(timesheetId, time, h.DB)
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
			employeeId, err := getEmployeeId(w, r)
			if err != nil {
				log.Println(err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			var requestParams struct {
				JobId         int    `json:"job_id"`
				WeekStartDate string `json:"week_start_date"`
			}

			err = decodeJSON(&requestParams, r)
			if err != nil {
				h.Logger.Errorf("Error decoding JSON: ", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			outTimesheetWeek, err := services.InitTimesheetWeek(
				employeeId,
				requestParams.JobId,
				requestParams.WeekStartDate,
				h.DB)
			if err != nil {
				log.Println("Error initializing timesheet week: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			responseJSON(w, outTimesheetWeek, h.Logger)
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

			timesheetWeeks, err := services.GetTimesheetWeekByEmployee(employeeId, h.DB)
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
			requestVals, err := parseRequestValues([]string{"id"}, r)
			if err != nil {
				log.Println("Error parsing query params: ", err)
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			timesheetWeek, err := services.DeleteTimesheetWeek(requestVals[0], h.DB)
			if err != nil {
				log.Println("Error deleting timesheet_week from db: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			responseJSON(w, timesheetWeek, h.Logger)
		},
	)
}
