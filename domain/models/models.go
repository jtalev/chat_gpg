package domain

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/nyaruka/phonenumbers"
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

type EmployeeAuth struct {
	AuthId       int    `json:"auth_id"`
	EmployeeId   string `json:"employee_id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type EmployeeErrors struct {
	EmployeeIdErr  string
	FirstNameErr   string
	LastNameErr    string
	EmailErr       string
	PhoneNumberErr string
	IsAdminErr     string
	IsSuccessful   bool
}

type EmployeeAuthErrors struct {
	UsernameErr  string
	PasswordErr  string
	IsSuccessful bool
}

func (e *Employee) Validate() EmployeeErrors {
	errors := EmployeeErrors{}
	errors.IsSuccessful = true
	errors = e.validateEmployeeId(errors)
	errors = e.validateFirstName(errors)
	errors = e.validateLastName(errors)
	errors = e.validateEmail(errors)
	errors = e.validatePhoneNumber(errors)
	errors = e.validateIsAdmin(errors)
	return errors
}

func (e *Employee) validateEmployeeId(errors EmployeeErrors) EmployeeErrors {
	if len(e.EmployeeId) <= 0 {
		errors.EmployeeIdErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	_, err := strconv.Atoi(e.EmployeeId)
	if err != nil {
		errors.EmployeeIdErr = "*must only contain numerical characters"
		errors.IsSuccessful = false
		return errors
	}
	if len(e.EmployeeId) != 7 {
		errors.EmployeeIdErr = "*must be 7 characters long"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (e *Employee) validateFirstName(errors EmployeeErrors) EmployeeErrors {
	if len(e.FirstName) <= 0 {
		errors.FirstNameErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (e *Employee) validateLastName(errors EmployeeErrors) EmployeeErrors {
	if len(e.LastName) <= 0 {
		errors.LastNameErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (e *Employee) validateEmail(errors EmployeeErrors) EmployeeErrors {
	if len(e.Email) <= 0 {
		errors.EmailErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	err := checkmail.ValidateFormat(e.Email)
	if err != nil {
		errors.EmailErr = "*invalid format"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (e *Employee) validatePhoneNumber(errors EmployeeErrors) EmployeeErrors {
	if len(e.PhoneNumber) <= 0 {
		errors.PhoneNumberErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	num, err := phonenumbers.Parse(e.PhoneNumber, "AU")
	if err != nil {
		errors.PhoneNumberErr = "*contains invalid characters"
		errors.IsSuccessful = false
		return errors
	}
	if !phonenumbers.IsValidNumber(num) {
		errors.PhoneNumberErr = "*invalid phone number"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (e *Employee) validateIsAdmin(errors EmployeeErrors) EmployeeErrors {
	if e.IsAdmin != true && e.IsAdmin != false {
		errors.IsAdminErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (e *EmployeeAuth) Validate() EmployeeAuthErrors {
	errors := EmployeeAuthErrors{}
	errors.IsSuccessful = true
	errors = e.validateUsername(errors)
	errors = e.validatePassword(errors)
	return errors
}

func (e *EmployeeAuth) validateUsername(errors EmployeeAuthErrors) EmployeeAuthErrors {
	if len(e.Username) <= 0 {
		errors.UsernameErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (e *EmployeeAuth) validatePassword(errors EmployeeAuthErrors) EmployeeAuthErrors {
	if len(e.PasswordHash) <= 0 {
		errors.PasswordErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	if len(e.PasswordHash) >= 0 && len(e.PasswordHash) <= 7 {
		errors.PasswordErr = "*password must be at least 8 characters long"
		errors.IsSuccessful = false
		return errors
	}
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

type LeaveErrors struct {
	DateErr      string
	IsSuccessful bool
}

func (l *LeaveRequest) Validate(pastRequests []LeaveRequest) LeaveErrors {
	errors := LeaveErrors{
		IsSuccessful: true,
	}
	errors = l.validateDate(errors, pastRequests)
	return errors
}

func stringToDate(date string) (time.Time, error) {
	date = strings.TrimSpace(date)
	dateArr := strings.Split(date, "-")
	dateArrInt := make([]int, 3)

	for i := range dateArr {
		int, err := strconv.Atoi(dateArr[i])
		if err != nil {
			return time.Time{}, errors.New("Date value cannot be converted to integer")
		}
		dateArrInt[i] = int
	}
	month := time.Month(dateArrInt[1])
	if month < time.January || month > time.December {
		return time.Time{}, errors.New("Invalid month value")
	}

	dateObject := time.Date(dateArrInt[0], month, dateArrInt[2], 0, 0, 0, 0, time.Local)

	return dateObject, nil
}

func (l *LeaveRequest) validateDate(errors LeaveErrors, pastRequests []LeaveRequest) LeaveErrors {
	from, err := stringToDate(l.From)
	if err != nil {
		errors.DateErr = "*invalid start date format"
		errors.IsSuccessful = false
		return errors
	}

	to, err := stringToDate(l.To)
	if err != nil {
		errors.DateErr = "*invalid end date format"
		errors.IsSuccessful = false
		return errors
	}

	if from.Before(time.Now().AddDate(0, 0, -1)) {
		errors.DateErr = "*start date cannot be a past date"
		errors.IsSuccessful = false
		return errors
	}

	if to.Before(from) {
		errors.DateErr = "*end date cannot be before start date"
		errors.IsSuccessful = false
		return errors
	}

	for _, pastLr := range pastRequests {
		pastFrom, err := stringToDate(pastLr.From)
		if err != nil {
			errors.DateErr = "*internal server error"
			errors.IsSuccessful = false
			return errors
		}
		pastTo, err := stringToDate(pastLr.To)
		if err != nil {
			errors.DateErr = "*internal server error"
			errors.IsSuccessful = false
			return errors
		}
		if pastFrom.Before(from) && pastTo.After(to) {
			errors.DateErr = "*dates overlap with previously submitted request"
			errors.IsSuccessful = false
			return errors
		}
		if pastFrom.After(from.AddDate(0, 0, -1)) && pastFrom.Before(to.AddDate(0, 0, 1)) {
			errors.DateErr = "*dates overlap with previously submitted request"
			errors.IsSuccessful = false
			return errors
		}
		if pastTo.After(from.AddDate(0, 0, -1)) && pastTo.Before(to.AddDate(0, 0, 1)) {
			errors.DateErr = "*dates overlap with previously submitted request"
			errors.IsSuccessful = false
			return errors
		}
	}

	return errors
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
	errors = j.validateName(errors)
	errors = j.validateNumber(errors)
	errors = j.validatePostCode(errors)
	return errors
}

func (j *Job) validateName(errors JobErrors) JobErrors {
	if len(j.Name) <= 0 {
		errors.IsSuccessful = false
		errors.NameErr = "*required"
		return errors
	}
	return errors
}

func (j *Job) validateNumber(errors JobErrors) JobErrors {
	if j.Number == 0 {
		return errors
	}
	if strconv.Itoa(j.Number) == "0" {
		errors.IsSuccessful = false
		errors.NumberErr = "*required"
		return errors
	}
	if strconv.Itoa(j.Number) < "0" {
		errors.IsSuccessful = false
		errors.NumberErr = "*invalid"
		return errors
	}
	return errors
}

func (j *Job) validateAddress(errors JobErrors) JobErrors {
	if len(j.Address) <= 0 {
		errors.IsSuccessful = false
		errors.AddressErr = "*required"
		return errors
	}
	return errors
}

func (j *Job) validateSuburb(errors JobErrors) JobErrors {
	if len(j.Suburb) <= 0 {
		errors.IsSuccessful = false
		errors.SuburbErr = "*required"
		return errors
	}
	return errors
}

func (j *Job) validatePostCode(errors JobErrors) JobErrors {
	if len(j.PostCode) == 0 {
		return errors
	}
	if j.PostCode == "n/a" {
		return errors
	}
	_, err := strconv.Atoi(j.PostCode)
	if err != nil {
		errors.IsSuccessful = false
		errors.PostCodeErr = "*numbers only"
		return errors
	}
	if len(j.PostCode) != 4 {
		errors.IsSuccessful = false
		errors.PostCodeErr = "*must be 4 characters long"
		return errors
	}

	return errors
}

func (j *Job) validateCity(errors JobErrors) JobErrors {
	if len(j.City) <= 0 {
		errors.IsSuccessful = false
		errors.CityErr = "*required"
		return errors
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
	log.Println(env)
	var dbPath string
	if env == "development" {
		dbPath = filepath.Join(rootPath, "infrastructure", "db", devFile)
	} else if env == "production" {
		dbPath = filepath.Join(rootPath, "infrastructure", "db", prodFile)
	}

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
