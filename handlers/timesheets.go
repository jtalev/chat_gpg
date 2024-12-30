package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jtalev/chat_gpg/models"
	"github.com/jtalev/chat_gpg/repository"
)

type TimesheetWeek struct {
	JobId      int
	Job        string
	Timesheets []models.Timesheet
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
			job, err := repository.GetJobById(timesheet.JobId, db)
			if err != nil {
				return nil, err
			}
			jobStr := fmt.Sprintf("%v, %v %v, %v", job.Name, job.Number, job.Address, job.Suburb)
			timesheetWeek := TimesheetWeek{
				JobId: timesheet.JobId,
				Job:   jobStr,
			}
			arr = append(arr, timesheetWeek)
			jobId = timesheet.JobId
		}
	}

	nullTimesheet := models.Timesheet{
		Hours:   0,
		Minutes: 00,
	}

	for i := range arr {
		arr[i].Timesheets = make([]models.Timesheet, 7)
		for j := range arr[i].Timesheets {
			arr[i].Timesheets[j] = nullTimesheet
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
					arr[i].Timesheets[0] = timesheet
				case "Thursday":
					arr[i].Timesheets[1] = timesheet
				case "Friday":
					arr[i].Timesheets[2] = timesheet
				case "Saturday":
					arr[i].Timesheets[3] = timesheet
				case "Sunday":
					arr[i].Timesheets[4] = timesheet
				case "Monday":
					arr[i].Timesheets[5] = timesheet
				case "Tuesday":
					arr[i].Timesheets[6] = timesheet
				default:
					return nil, errors.New("Unexpected value for day.String()")
				}
			}
		}
	}

	return arr, nil
}

func (h *Handler) RenderTimesheetByWeek() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("rendering timesheet row")
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

			arrow := r.FormValue("arrow")
			wedDateStr := r.FormValue("wedDate")
			wedDate, err := strconv.Atoi(wedDateStr)
			if err != nil {
				h.Logger.Error("Error converting date to int")
				http.Error(w, "Error bad request", http.StatusBadRequest)
				return
			}
			monthStr := r.FormValue("month")
			month, err := strconv.Atoi(monthStr)
			if err != nil {
				h.Logger.Error("Error converting month to int")
				http.Error(w, "Error bad request", http.StatusBadRequest)
				return
			}
			timeMonth := time.Month(month)
			yearStr := r.FormValue("year")
			year, err := strconv.Atoi(yearStr)
			if err != nil {
				h.Logger.Error("Error converting year to int")
				http.Error(w, "Error bad request", http.StatusBadRequest)
				return
			}

			date := time.Date(year, timeMonth, wedDate, 0, 0, 0, 0, time.Local)
			if arrow == "left" {
				date = date.AddDate(0, 0, -7)
			} else if arrow == "right" {
				date = date.AddDate(0, 0, 7)
			}

			year, timeMonth, wedDate = date.Date()
			weekStart := fmt.Sprintf("%v-%v-%v", year, int(timeMonth), wedDate)

			data, err := mapTimesheets(employeeId, weekStart, h.DB)
			if err != nil {
				h.Logger.Errorf("Error mapping timesheet week data: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			fmt.Println(data)

			timesheetRowPath := filepath.Join("..", "ui", "templates", "timesheetRow.html")
			tmpl, err := template.ParseFiles(timesheetRowPath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = tmpl.ExecuteTemplate(w, "timesheetRow", data)
			if err != nil {
				h.Logger.Errorf("Error executing template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}
