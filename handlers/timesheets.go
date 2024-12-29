package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/jtalev/chat_gpg/models"
	"github.com/jtalev/chat_gpg/repository"
)

type TimesheetWeek struct {
	JobId      int
	Job        string
	Timesheets []map[string]models.Timesheet
}

func mapTimesheets(employeeId, weekStart string, db *sql.DB) ([]TimesheetWeek, error) {
	timesheets, err := repository.GetTimesheetsByWeekStart(employeeId, weekStart, db)
	if err != nil {
		return nil, err
	}

	jobId := 0
	var arr []TimesheetWeek

	for _, timesheet := range timesheets {
		if timesheet.JobId != jobId {
			newJob, err := repository.GetJobById(timesheet.JobId, db)
			if err != nil {
				return nil, err
			}
			job := newJob.Name + ", " + string(newJob.Number) + " " + newJob.Address + ", " + newJob.Suburb
			timesheetWeek := TimesheetWeek{
				JobId: timesheet.JobId,
				Job:   job,
			}
			arr = append(arr, timesheetWeek)
			jobId = timesheet.JobId
		}
	}

	for i := range arr {
		arr[i].Timesheets = make([]map[string]models.Timesheet, 7)
		for j := range arr[i].Timesheets {
			arr[i].Timesheets[j] = make(map[string]models.Timesheet)
		}
	}

	for _, timesheet := range timesheets {
		for i := range arr {
			if timesheet.JobId == arr[i].JobId {
				dateStr := timesheet.Date
				date, err := stringToDate(dateStr)
				if err != nil {
					return nil, err
				}
				day := date.Weekday()

				switch day.String() {
				case "Wednesday":
					arr[i].Timesheets[0][day.String()] = timesheet
				case "Thursday":
					arr[i].Timesheets[1][day.String()] = timesheet
				case "Friday":
					arr[i].Timesheets[2][day.String()] = timesheet
				case "Saturday":
					arr[i].Timesheets[3][day.String()] = timesheet
				case "Sunday":
					arr[i].Timesheets[4][day.String()] = timesheet
				case "Monday":
					arr[i].Timesheets[5][day.String()] = timesheet
				case "Tuesday":
					arr[i].Timesheets[6][day.String()] = timesheet
				default:
					return nil, errors.New("Unexpected value for day.String()")
				}
			}
		}
	}

	for _, timesheetWeek := range arr {
		fmt.Println(timesheetWeek)
	}

	return arr, nil
}

func (h *Handler) GetTimesheetsByWeekStart() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				h.Logger.Errorf("Error parsing form: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			employeeId, ok := r.Context().Value("employee_id").(string)
			if !ok {
				h.Logger.Error("Error getting employee_id from context")
				http.Error(w, "Error getting value context", http.StatusNotFound)
				return
			}

			weekStart := r.FormValue("weekStart")

			data, err := mapTimesheets(employeeId, weekStart, h.DB)
			if err != nil {
				h.Logger.Errorf("Error mapping timesheet week data: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			responseJson(w, data, h.Logger)
		},
	)
}
