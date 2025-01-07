package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jtalev/chat_gpg/models"
	"github.com/jtalev/chat_gpg/repository"
	"go.uber.org/zap"
)

func getEmployeeId(w http.ResponseWriter, r *http.Request) (string, error) {
	employeeId, ok := r.Context().Value("employee_id").(string)
	if !ok {
		return "", errors.New("Error retrieving employee_id")
	}
	return employeeId, nil
}

func responseJSON(w http.ResponseWriter, data any, sugar *zap.SugaredLogger) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		sugar.Errorf("Error encoding leave requests: %v", err)
		http.Error(w, "failed to fetch leave requests", http.StatusInternalServerError)
		return
	}
}

func executePartialTemplate(filepath string, data interface{}, w http.ResponseWriter) error {
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		return err
	}
	err = tmpl.ExecuteTemplate(w, "timesheetTable", data)
	if err != nil {
		return err
	}
	return nil
}

func renderTemplate(
	w http.ResponseWriter,
	r *http.Request,
	component, title string,
	componentData interface{},
) {
	layoutPath := filepath.Join("..", "ui", "layouts", "layout.html")
	navPath := filepath.Join("..", "ui", "templates", "nav.html")
	dashboardPath := filepath.Join("..", "ui", "views", "dashboard.html")
	jobsPath := filepath.Join("..", "ui", "views", "jobs.html")
	timesheetsPath := filepath.Join("..", "ui", "views", "timesheets.html")
	timesheetTablePath := filepath.Join("..", "ui", "templates", "timesheetTable.html")
	leavePath := filepath.Join("..", "ui", "views", "leave.html")
	adminPath := filepath.Join("..", "ui", "views", "admin.html")
	accountPath := filepath.Join("..", "ui", "views", "account.html")

	isAdmin, ok := r.Context().Value("is_admin").(bool)
	if !ok {
		http.Error(w, "unable to retrieve is_admin", http.StatusUnauthorized)
		return
	}

	data := struct {
		Title     string
		Component string
		IsAdmin   bool
		Data      interface{}
	}{
		Title:     title,
		Component: component,
		IsAdmin:   isAdmin,
		Data:      componentData,
	}

	tmpl, err := template.ParseFiles(
		layoutPath,
		navPath,
		dashboardPath,
		jobsPath,
		timesheetsPath,
		timesheetTablePath,
		leavePath,
		adminPath,
		accountPath,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

type Handler struct {
	DB     *sql.DB
	Store  *sessions.CookieStore
	Logger *zap.SugaredLogger
}

func (h *Handler) ServeLoginView(w http.ResponseWriter, r *http.Request) {
	login_path := filepath.Join("..", "ui", "views", "login.html")
	tmpl := template.Must(template.ParseFiles(login_path))
	tmpl.Execute(w, nil)
}

func (h *Handler) ServeErrorView(w http.ResponseWriter, r *http.Request) {
	errorPath := filepath.Join("..", "ui", "views", "error.html")
	tmpl := template.Must(template.ParseFiles(errorPath))
	tmpl.Execute(w, nil)
}

func (h *Handler) ServeAccountView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			component := "account"
			title := "Account - GPG"
			renderTemplate(w, r, component, title, nil)
		},
	)
}

func (h *Handler) ServeAdminView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data := getAdminData()
			component := "admin"
			title := "Admin - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}

func (h *Handler) ServeDashboardView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data := getDashboardData()
			component := "dashboard"
			title := "Dashboard - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}

func (h *Handler) ServeJobsView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data := getJobsData()
			component := "jobs"
			title := "Jobs - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}

type LeaveViewData struct {
	Employee      models.Employee
	LeaveRequests []models.LeaveRequest
	DateError     string
}

func (h *Handler) ServeLeaveView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			employeeId, ok := r.Context().Value("employee_id").(string)
			if !ok {
				h.Logger.Warn("Error getting employee_id from context")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			employee, err := repository.GetEmployeeByEmployeeId(employeeId, h.DB)
			if err != nil {
				h.Logger.Errorf("Error getting employee data: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			leaveRequests, err := repository.GetLeaveRequestsByEmployee(employeeId, h.DB)
			if err != nil {
				h.Logger.Errorf("Error getting leave page data: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			data := LeaveViewData{
				Employee:      employee,
				LeaveRequests: leaveRequests,
				DateError:     "",
			}

			component := "leave"
			title := "Leave - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}

type TimesheetViewData struct {
	Jobs              []models.Job
	InitialTimesheets []TimesheetWeek
	WedDate           int
	MonthInt          int
	MonthStr          string
	Year              int
	PrevWeekDates     []int
}

func wednesdayDate() (int, time.Month, int, int) {
	now := time.Now()
	weekday := now.Weekday()
	year, month, day := now.Date()
	date := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

	switch weekday {
	case time.Sunday:
		date = date.AddDate(0, 0, -4)
	case time.Monday:
		date = date.AddDate(0, 0, -5)
	case time.Tuesday:
		date = date.AddDate(0, 0, -6)
	case time.Wednesday:
		date = date.AddDate(0, 0, 0)
	case time.Thursday:
		date = date.AddDate(0, 0, -1)
	case time.Friday:
		date = date.AddDate(0, 0, -2)
	case time.Saturday:
		date = date.AddDate(0, 0, -3)
	default:
		fmt.Println("no dates added")
	}

	year, month, day = date.Date()

	return year, month, int(month), day
}

func (h *Handler) ServeTimesheetsView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			jobs, err := repository.GetJobs(h.DB)
			if err != nil {
				h.Logger.Errorf("Error querying job database: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			employeeId, ok := r.Context().Value("employee_id").(string)
			if !ok {
				h.Logger.Warn("Error getting employee_id from context")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			year, month, monthInt, day := wednesdayDate()
			weekStart := fmt.Sprintf("%v-%v-%v", year, monthInt, day)

			timesheets, err := mapTimesheets(employeeId, weekStart, h.DB)
			if err != nil {
				h.Logger.Errorf("Error getting initial timesheet data: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			prevWeekDates := prevWeekDates(time.Date(year, month, day, 0, 0, 0, 0, time.Local))

			data := TimesheetViewData{
				Jobs:              jobs,
				InitialTimesheets: timesheets,
				WedDate:           day,
				MonthInt:          monthInt,
				MonthStr:          month.String(),
				Year:              year,
				PrevWeekDates:     prevWeekDates,
			}
			component := "timesheets"
			title := "Timesheets - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}
