package domain

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"go.uber.org/zap"
)

type ValidationResult struct {
	Key     string
	IsValid bool
	Msg     string
}

type Employee struct {
	ID          int    `json:"id"`
	EmployeeId  string `json:"employee_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	IsAdmin     bool   `json:"is_admin"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type EmployeeErrors struct {
	EmployeeIdErr  string
	FirstNameErr   string
	LastNameErr    string
	EmailErr       string
	PhoneNumberErr string
	IsSuccessful   bool
}

type EmployeeAuthErrors struct {
	UsernameErr  string
	PasswordErr  string
	IsSuccessful bool
}

type EmployeeAuth struct {
	AuthId       int    `json:"auth_id"`
	EmployeeId   string `json:"employee_id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func (e *Employee) Validate() EmployeeErrors {
	errors := EmployeeErrors{}
	errors.IsSuccessful = true
	return errors
}

type LeaveRequest struct {
	RequestId  int    `json:"request_id"`
	EmployeeId string `json:"employee_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Type       string `json:"leave_type"`
	From       string `json:"from_date"`
	To         string `json:"to_date"`
	Note       string `json:"note"`
	IsPending  bool   `json:"is_pending"`
	IsApproved bool   `json:"is_approved"`
}

type Job struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Number     int    `json:"number"`
	Address    string `json:"address"`
	Suburb     string `json:"suburb"`
	PostCode   string `json:"post_code"`
	City       string `json:"city"`
	IsComplete bool   `json:"is_available"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type JobErrors struct {
	NameErr      string
	NumberErr    string
	AddressErr   string
	SuburbErr    string
	PostCodeErr  string
	CityErr      string
	IsSuccessful bool
}

func (j *Job) Validate() JobErrors {
	errors := JobErrors{}
	errors.IsSuccessful = true
	errors = j.ValidatePostCode(errors)
	return errors
}

func (j *Job) ValidatePostCode(errors JobErrors) JobErrors {
	_, err := strconv.Atoi(j.PostCode)
	if err != nil {
		errors.IsSuccessful = false
		errors.PostCodeErr = "Numbers only"
	}
	log.Println(len(j.PostCode))
	if len(j.PostCode) != 4 {
		errors.IsSuccessful = false
		errors.PostCodeErr = "Must be 4 characters long"
	}

	return errors
}

type TimesheetWeek struct {
	TimesheetWeekId int    `json:"timesheet_week_id"`
	EmployeeId      string `json:"employee_id"`
	JobId           int    `json:"job_id"`
	WedTimesheetId  int    `json:"wed_timesheet_id"`
	ThuTimesheetId  int    `json:"thu_timesheet_id"`
	FriTimesheetId  int    `json:"fri_timesheet_id"`
	SatTimesheetId  int    `json:"sat_timesheet_id"`
	SunTimesheetId  int    `json:"sun_timesheet_id"`
	MonTimesheetId  int    `json:"mon_timesheet_id"`
	TueTimesheetId  int    `json:"tue_timesheet_id"`
	WeekStartDate   string `json:"week_start_date"`
	Total           int    `json:"total"`
	CreatedAt       string `json:"created_at"`
	ModifiedAt      string `json:"modified_at"`
}

type Timesheet struct {
	TimesheetId     int    `json:"timesheet_id"`
	TimesheetWeekId int    `json:"timesheet_week_id"`
	TimesheetDate   string `json:"timesheet_date"`
	Day             string `json:"day"`
	Hours           int    `json:"hours"`
	Minutes         int    `json:"minutes"`
	CreatedAt       string `json:"created_at"`
	ModifiedAt      string `json:"modified_at"`
}

func InitDb(rootPath string, sugar *zap.SugaredLogger) *sql.DB {
	devFile := "dev.db"
	prodFile := "prod.db"
	env := os.Getenv("ENV")
	var dbPath string
	if env == "development" {
		dbPath = filepath.Join(rootPath, "infrastructure", "db", devFile)
	} else if env == "production" {
		dbPath = filepath.Join(rootPath, "infrastructure", "db", prodFile)
	}

	log.Println(dbPath)

	if dbPath == "" {
		sugar.Error("Error obtaining db path")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		sugar.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		sugar.Error("DB connection not open:", err)
	} else {
		sugar.Info("DB connection is open and healthy")
	}

	return db
}
