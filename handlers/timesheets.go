package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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
	isJobAdded := false

	for _, timesheet := range timesheets {
		if timesheet.JobId != jobId {
			for i := range arr {
				if arr[i].JobId == timesheet.JobId {
					isJobAdded = true
				}
			}
			if isJobAdded {
				isJobAdded = false
				continue
			}
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

type Data struct {
	Data TimesheetTemplateData
}

type TimesheetTemplateData struct {
	Jobs              []models.Job
	InitialTimesheets []TimesheetWeek
	WedDate           int
	MonthInt          int
	MonthStr          time.Month
	Year              int
	PrevWeekDates     []int
}

func prevWeekDates(date time.Time) []int {
	dates := make([]int, 7)
	for i := range dates {
		_, _, day := date.Date()
		dates[i] = day
		date = date.AddDate(0, 0, 1)
	}
	return dates
}

func (h *Handler) RenderTimesheetByWeek() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			jobs, err := repository.GetJobs(h.DB)
			if err != nil {
				h.Logger.Errorf("Error querying job database: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			err = r.ParseForm()
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

			initialTimesheets, err := mapTimesheets(employeeId, weekStart, h.DB)
			if err != nil {
				h.Logger.Errorf("Error mapping timesheet week data: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			prevWeekDates := prevWeekDates(date)

			var data Data
			data.Data.Jobs = jobs
			data.Data.InitialTimesheets = initialTimesheets
			data.Data.WedDate = wedDate
			data.Data.MonthInt = int(timeMonth)
			data.Data.MonthStr = timeMonth
			data.Data.Year = year
			data.Data.PrevWeekDates = prevWeekDates

			timesheetTablePath := filepath.Join("..", "ui", "templates", "timesheetTable.html")
			tmpl, err := template.ParseFiles(timesheetTablePath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = tmpl.ExecuteTemplate(w, "timesheetTable", data)
			if err != nil {
				h.Logger.Errorf("Error executing template: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) GetTimesheetById() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				h.Logger.Errorf("Error parsing form: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			idStr := r.FormValue("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				h.Logger.Errorf("Invalid form value: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			timesheet, err := repository.GetTimesheetById(id, h.DB)
			if err != nil {
				h.Logger.Errorf("Error querying timesheet from db: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			responseJson(w, timesheet, h.Logger)
		},
	)
}

func strToInt(str []string) ([]int, error) {
	res := make([]int, len(str))

	for i := range str {
		result, err := strconv.Atoi(str[i])
		if err != nil {
			return nil, err
		}
		res[i] = result
	}
	return res, nil
}

func monthFromString(str string) (time.Month, error) {
	monthMap := map[string]time.Month{
		"January":   time.January,
		"February":  time.February,
		"March":     time.March,
		"April":     time.April,
		"May":       time.May,
		"June":      time.June,
		"July":      time.July,
		"August":    time.August,
		"September": time.September,
		"October":   time.October,
		"November":  time.November,
		"December":  time.December,
	}

	str = strings.Title(strings.ToLower(str))
	if month, ok := monthMap[str]; ok {
		return month, nil
	}
	return 0, errors.New("input string doesn't match map")
}

func (h *Handler) PostTimesheetsAll() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			employeeId, ok := r.Context().Value("employee_id").(string)
			if !ok {
				h.Logger.Error("Error getting employee_id from context")
				http.Error(w, "Error getting value context", http.StatusNotFound)
				return
			}

			var payload struct {
				PayloadTimesheets []struct {
					JobId     string `json:"job"`
					Time      string `json:"time"`
					WeekStart struct {
						Date  string `json:"date"`
						Month string `json:"month"`
						Year  string `json:"year"`
					} `json:"weekStart"`
					Date string `json:"date"`
				} `json:"timesheets"`
			}

			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				h.Logger.Errorf("Error decoding JSON: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			var timesheetsToPost []models.Timesheet
			for _, timesheet := range payload.PayloadTimesheets {
				var ts models.Timesheet

				dateArrStr := make([]string, 3)
				dateArrStr[0] = timesheet.WeekStart.Date
				dateArrStr[1] = timesheet.WeekStart.Month
				dateArrStr[2] = timesheet.WeekStart.Year
				month, err := monthFromString(dateArrStr[1])
				if err != nil {
					h.Logger.Errorf("Error converting month str to time.Month: %v", err)
					http.Error(w, "Bad request", http.StatusBadRequest)
					return
				}
				weekStart := fmt.Sprintf("%v-%v-%v", dateArrStr[2], int(month), dateArrStr[0])

				dateArrStr[1] = "0"
				dateArrInt, err := strToInt(dateArrStr)
				if err != nil {
					h.Logger.Error("Payload string can't be converted to int: %v", err)
					http.Error(w, "Bad request", http.StatusBadRequest)
					return
				}

				date := time.Date(dateArrInt[2], month, dateArrInt[0], 0, 0, 0, 0, time.Local)
				tsDateStr := make([]string, 1)
				tsDateStr[0] = timesheet.Date
				tsDateInt, err := strToInt(tsDateStr)
				if err != nil {
					h.Logger.Error("Payload string can't be converted to int: %v", err)
					http.Error(w, "Bad request", http.StatusBadRequest)
					return
				}
				tsDate := ""
				for tsDate == "" {
					if date.Day() == tsDateInt[0] {
						tsDate = fmt.Sprintf("%v-%v-%v", date.Year(), int(date.Month()), date.Day())
					} else {
						date = date.AddDate(0, 0, 1)
					}
				}

				var timeStr []string
				if strings.Contains(timesheet.Time, ":") {
					timeStr = strings.Split(timesheet.Time, ":")
					if timeStr[0] == "" {
						timeStr[0] = "0"
					} else if timeStr[1] == "" {
						timeStr[1] = "0"
					}
				} else {
					timeStr = append(timeStr, timesheet.Time)
					timeStr = append(timeStr, "0")
				}
				timeInt, err := strToInt(timeStr)
				if err != nil {
					h.Logger.Error("Payload string can't be converted to int: %v", err)
					http.Error(w, "Bad request", http.StatusBadRequest)
					return
				}

				ts.EmployeeId = employeeId
				ts.JobId, err = strconv.Atoi(timesheet.JobId)
				if err != nil {
					h.Logger.Error("Payload string can't be converted to int: %v", err)
					http.Error(w, "Bad request", http.StatusBadRequest)
					return
				}
				ts.WeekStart = weekStart
				ts.Date = tsDate
				ts.Hours = timeInt[0]
				ts.Minutes = timeInt[1]

				timesheetsToPost = append(timesheetsToPost, ts)
			}

			for _, timesheet := range timesheetsToPost {
				_, err := repository.PostTimesheet(timesheet, h.DB)
				if err != nil {
					h.Logger.Errorf("Error posting timesheet to db: %v", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
			}
		},
	)
}

func (h *Handler) PutTimesheetsAll() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var payload struct {
				PayloadTimesheets []struct {
					ID   string `json:"id"`
					Time string `json:"time"`
				} `json:"timesheets"`
			}

			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				h.Logger.Errorf("Error decoding JSON: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			var timesheetsToUpdate []models.Timesheet
			for _, timesheet := range payload.PayloadTimesheets {
				str := strings.Split(timesheet.Time, ":")
				str = append([]string{timesheet.ID}, str...)
				intArr, err := strToInt(str)
				if err != nil {
					h.Logger.Error("Payload string can't be converted to int: %v", err)
					http.Error(w, "Bad request", http.StatusBadRequest)
					return
				}
				ts := models.Timesheet{
					ID:      intArr[0],
					Hours:   intArr[1],
					Minutes: intArr[2],
				}
				timesheetsToUpdate = append(timesheetsToUpdate, ts)
			}

			for _, timesheet := range timesheetsToUpdate {
				_, err := repository.PutTimesheet(timesheet, h.DB)
				if err != nil {
					h.Logger.Errorf("Error updating timesheet: %v", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
			}
		},
	)
}

func (h *Handler) PutTimesheet() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var data struct {
				ID   int    `json:"id"`
				Time string `json:"time"`
			}

			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				h.Logger.Errorf("Error decoding JSON: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			str := strings.Split(data.Time, ":")
			intArr, err := strToInt(str)
			if err != nil {
				h.Logger.Error("Payload string can't be converted to int: %v", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			hours := intArr[0]
			minutes := intArr[1]

			timesheet := models.Timesheet{
				ID:      data.ID,
				Hours:   hours,
				Minutes: minutes,
			}

			timesheet, err = repository.PutTimesheet(timesheet, h.DB)
			if err != nil {
				h.Logger.Errorf("Error updating timesheet: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			responseJson(w, timesheet, h.Logger)
		},
	)
}
