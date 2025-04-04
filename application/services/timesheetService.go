package application

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	domain "github.com/jtalev/chat_gpg/domain/models"
	infrastructure "github.com/jtalev/chat_gpg/infrastructure/repository"
)

var nilTimesheet = domain.Timesheet{
	TimesheetId:     0,
	TimesheetWeekId: 0,
	TimesheetDate:   "",
	Day:             "",
	Hours:           0,
	Minutes:         0,
	CreatedAt:       "",
	ModifiedAt:      "",
}

var days = []string{
	"wed",
	"thu",
	"fri",
	"sat",
	"sun",
	"mon",
	"tue",
}

func splitTime(inTime string) ([]string, error) {
	var outTimeArr []string
	containsColon := strings.Contains(inTime, ":")
	if !containsColon {
		outTimeArr = append(outTimeArr, inTime, "0")
	} else {
		outTimeArr = strings.Split(inTime, ":")
	}
	return outTimeArr, nil
}

func strToInt(inArr []string) ([]int, error) {
	outArr := make([]int, len(inArr))
	for i := range inArr {
		integer, err := strconv.Atoi(inArr[i])
		if err != nil {
			return nil, err
		}
		outArr[i] = integer
	}
	return outArr, nil
}

// incoming dates are in format yyyy-mm-dd
func DateStrToDate(inDate string) (time.Time, error) {
	containsHyphen := strings.Contains(inDate, "-")
	if !containsHyphen {
		return time.Time{}, errors.New("Date string must be in format: yyyy-mm-dd")
	}

	dateStrArr := strings.Split(inDate, "-")
	dateIntArr, err := strToInt(dateStrArr)
	if err != nil {
		return time.Time{}, err
	}
	year, month, day := dateIntArr[0], dateIntArr[1], dateIntArr[2]
	monthMonth := time.Month(month)
	outDate := time.Date(year, monthMonth, day, 0, 0, 0, 0, time.Local)
	return outDate, nil
}

func GetAllTimesheets(db *sql.DB) ([]domain.Timesheet, error) {
	outTimesheets, err := infrastructure.GetTimesheets(db)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return outTimesheets, nil
}

func GetTimesheetById(inTimesheetId string, db *sql.DB) (domain.Timesheet, error) {
	inId, err := strconv.Atoi(inTimesheetId)
	if err != nil {
		log.Println(err)
		return domain.Timesheet{}, err
	}

	outTimesheet, err := infrastructure.GetTimesheetById(inId, db)
	if err != nil {
		log.Println(err)
		return domain.Timesheet{}, err
	}

	return outTimesheet, nil
}

func PostTimesheet(inTimesheet domain.Timesheet, db *sql.DB) (domain.Timesheet, error) {
	if inTimesheet.TimesheetId == 0 {
		inTimesheet = nilTimesheet
	}

	// TODO: validate timesheet

	outTimesheet, err := infrastructure.PostTimesheet(inTimesheet, db)
	if err != nil {
		log.Println(err)
		return outTimesheet, err
	}

	return outTimesheet, err
}

func PutTimesheet(id, time string, db *sql.DB) (domain.Timesheet, error) {
	idConv, err := strconv.Atoi(id)
	if err != nil {
		return domain.Timesheet{}, err
	}
	parsedTime, err := splitTime(time)
	if err != nil {
		return domain.Timesheet{}, nil
	}

	parsedTimeInt, err := strToInt(parsedTime)
	if err != nil {
		return domain.Timesheet{}, err
	}

	inTimesheet := domain.Timesheet{
		TimesheetId: idConv,
		Hours:       parsedTimeInt[0],
		Minutes:     parsedTimeInt[1],
	}

	// TODO: validate timesheet

	outTimesheet, err := infrastructure.PutTimesheet(inTimesheet, db)
	if err != nil {
		log.Println(err)
		return outTimesheet, err
	}

	return outTimesheet, err
}

func postNilTimesheets(inTimesheetWeek *domain.TimesheetWeek, weekStartDate string, db *sql.DB) ([]domain.Timesheet, error) {
	var outTimesheets = make([]domain.Timesheet, 7)
	timesheetDate := weekStartDate
	date, err := DateStrToDate(weekStartDate)
	if err != nil {
		return nil, err
	}

	for i := range outTimesheets {
		outTimesheets[i].TimesheetWeekId = inTimesheetWeek.TimesheetWeekId
		outTimesheets[i].TimesheetDate = timesheetDate
		outTimesheets[i].Day = days[i]
		outTimesheet, err := infrastructure.PostTimesheet(outTimesheets[i], db)
		if err != nil {
			return nil, err
		}

		date = date.AddDate(0, 0, 1)
		timesheetDate = fmt.Sprintf("%v-%v-%v", date.Year(), int(date.Month()), date.Day())
		outTimesheets[i] = outTimesheet
	}

	return outTimesheets, nil
}

func putTimesheetWeek(inTimesheetWeek domain.TimesheetWeek, db *sql.DB) (domain.TimesheetWeek, error) {
	outTimesheetWeek, err := infrastructure.PutTimesheetWeek(inTimesheetWeek, db)
	if err != nil {
		return domain.TimesheetWeek{}, err
	}
	return outTimesheetWeek, nil
}

func InitTimesheetWeek(employeeId string, jobIdStr string, weekStartDate string, db *sql.DB) ([]TimesheetRow, error) {
	jobId, err := strconv.Atoi(jobIdStr)
	if err != nil {
		return nil, err
	}

	inTimesheetWeek := domain.TimesheetWeek{
		EmployeeId:    employeeId,
		JobId:         jobId,
		WeekStartDate: weekStartDate,
	}

	outTimesheetWeek, err := infrastructure.PostTimesheetWeek(inTimesheetWeek, db)
	if err != nil {
		return nil, nil
	}
	initialTimesheets, err := postNilTimesheets(&outTimesheetWeek, weekStartDate, db)
	outTimesheetWeek.WedTimesheetId = initialTimesheets[0].TimesheetId
	outTimesheetWeek.ThuTimesheetId = initialTimesheets[1].TimesheetId
	outTimesheetWeek.FriTimesheetId = initialTimesheets[2].TimesheetId
	outTimesheetWeek.SatTimesheetId = initialTimesheets[3].TimesheetId
	outTimesheetWeek.SunTimesheetId = initialTimesheets[4].TimesheetId
	outTimesheetWeek.MonTimesheetId = initialTimesheets[5].TimesheetId
	outTimesheetWeek.TueTimesheetId = initialTimesheets[6].TimesheetId

	outTimesheetWeek, err = infrastructure.PutTimesheetWeek(outTimesheetWeek, db)
	if err != nil {
		return nil, err
	}

	outTimesheetRows, err := MapTimesheetsToTimesheetWeek([]domain.TimesheetWeek{outTimesheetWeek}, db)
	if err != nil {
		return nil, err
	}

	return outTimesheetRows, nil
}

func GetTimesheetWeekByEmployee(employeeId string, db *sql.DB) ([]domain.TimesheetWeek, error) {

	// TODO: validate employeeId
	outTimesheetWeeks, err := infrastructure.GetTimesheetWeekByEmployee(employeeId, db)
	if err != nil {
		return nil, err
	}
	return outTimesheetWeeks, nil
}

func GetTimesheetWeekByEmployeeWeekStart(employeeId, weekStart string, db *sql.DB) ([]domain.TimesheetWeek, error) {
	// TODO: validate weekStart is in correct format
	outTimesheetWeeks, err := infrastructure.GetTimesheetWeekByEmployeeWeekStart(employeeId, weekStart, db)
	if err != nil {
		return nil, err
	}
	return outTimesheetWeeks, nil
}

func DeleteTimesheetWeek(id string, db *sql.DB) (domain.TimesheetWeek, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return domain.TimesheetWeek{}, err
	}
	outTimesheetWeek, err := infrastructure.DeleteTimesheetWeek(idInt, db)
	if err != nil {
		return outTimesheetWeek, err
	}
	return outTimesheetWeek, nil
}

// data needed for existingTimesheetRow.html
type TimesheetRow struct {
	JobName         string
	Timesheets      []domain.Timesheet
	TimesheetWeekId int
	Total           string
}

type TimesheetViewData struct {
	MonthStr      string
	Year          int //not sure if int or string yet
	Dates         []int
	TimesheetRows []TimesheetRow
	WeekStartDate string
}

func InitialTimesheetViewData(employeeId string, db *sql.DB) (TimesheetViewData, error) {
	outData := TimesheetViewData{}

	dates, err := currentWeekDates()
	if err != nil {
		log.Println("Error getting current week dates")
		return outData, err
	}

	year, month, day := WeekStartDate().Date()
	weekStart := fmt.Sprintf("%v-%v-%v", year, int(month), day)
	initialTimesheetWeeks, err := GetTimesheetWeekByEmployeeWeekStart(employeeId, weekStart, db)
	if err != nil {
		return outData, err
	}
	timesheetRows, err := MapTimesheetsToTimesheetWeek(initialTimesheetWeeks, db)

	outData.MonthStr = month.String()
	outData.Year = year
	outData.Dates = dates
	outData.TimesheetRows = timesheetRows
	outData.WeekStartDate = weekStart

	return outData, nil
}

func MapTimesheetsToTimesheetWeek(inTimesheetWeeks []domain.TimesheetWeek, db *sql.DB) ([]TimesheetRow, error) {
	outData := make([]TimesheetRow, len(inTimesheetWeeks))
	for i := range inTimesheetWeeks {
		job, err := GetJobById(inTimesheetWeeks[i].JobId, db)
		outData[i].TimesheetWeekId = inTimesheetWeeks[i].TimesheetWeekId
		if err != nil {
			return nil, err
		}
		jobName := fmt.Sprintf("%s", job.Name)
		if job.Number != 0 && job.Address != "n/a" {
			jobName = fmt.Sprintf("%s, %v %s", jobName, job.Number, job.Address)
		}
		if job.PostCode != "n/a" {
			jobName = fmt.Sprintf("%s, %s", jobName, job.PostCode)
		}
		outData[i].JobName = jobName

		wed, err := GetTimesheetById(strconv.Itoa(inTimesheetWeeks[i].WedTimesheetId), db)
		if err != nil {
			return nil, err
		}
		thu, err := GetTimesheetById(strconv.Itoa(inTimesheetWeeks[i].ThuTimesheetId), db)
		if err != nil {
			return nil, err
		}
		fri, err := GetTimesheetById(strconv.Itoa(inTimesheetWeeks[i].FriTimesheetId), db)
		if err != nil {
			return nil, err
		}
		sat, err := GetTimesheetById(strconv.Itoa(inTimesheetWeeks[i].SatTimesheetId), db)
		if err != nil {
			return nil, err
		}
		sun, err := GetTimesheetById(strconv.Itoa(inTimesheetWeeks[i].SunTimesheetId), db)
		if err != nil {
			return nil, err
		}
		mon, err := GetTimesheetById(strconv.Itoa(inTimesheetWeeks[i].MonTimesheetId), db)
		if err != nil {
			return nil, err
		}
		tue, err := GetTimesheetById(strconv.Itoa(inTimesheetWeeks[i].TueTimesheetId), db)
		if err != nil {
			return nil, err
		}
		outData[i].Timesheets = append(outData[i].Timesheets, wed, thu, fri, sat, sun, mon, tue)
	}
	return outData, nil
}

func WeekStartDate() time.Time {
	date := time.Now()
	switch date.Weekday().String() {
	case "Wednesday":
		break
	case "Thursday":
		date = date.AddDate(0, 0, -1)
	case "Friday":
		date = date.AddDate(0, 0, -2)
	case "Saturday":
		date = date.AddDate(0, 0, -3)
	case "Sunday":
		date = date.AddDate(0, 0, -4)
	case "Monday":
		date = date.AddDate(0, 0, -5)
	case "Tuesday":
		date = date.AddDate(0, 0, -6)
	}
	return date
}

func currentWeekDates() ([]int, error) {
	date := WeekStartDate()
	outDates := make([]int, 7)

	for i := range outDates {
		outDates[i] = date.Day()
		date = date.AddDate(0, 0, 1)
	}
	return outDates, nil
}

type Data struct {
	MonthStr      string
	Year          int
	Dates         []int
	TimesheetRows []TimesheetRow
	WeekStartDate string
}

type TimesheetData struct {
	Data Data
}

func TimesheetTableData(employeeId string, hxVals []string, db *sql.DB) (TimesheetData, error) {
	outData := TimesheetData{}

	arrow, weekStartDate := hxVals[0], hxVals[1]
	date, err := DateStrToDate(weekStartDate)
	if err != nil {
		log.Println(err)
		return TimesheetData{}, err
	}
	if arrow == "left" {
		date = date.AddDate(0, 0, -7)
	} else if arrow == "right" {
		date = date.AddDate(0, 0, 7)
	} else {
		return TimesheetData{}, err
	}
	monthStr := date.Month().String()
	year := date.Year()
	weekStartDate = fmt.Sprintf("%v-%v-%v", date.Year(), int(date.Month()), date.Day())
	dates := make([]int, 7)
	for i := range dates {
		dates[i] = date.Day()
		date = date.AddDate(0, 0, 1)
	}
	timesheetWeeks, err := infrastructure.GetTimesheetWeekByEmployeeWeekStart(employeeId, weekStartDate, db)
	if err != nil {
		return outData, err
	}
	timesheetRows, err := MapTimesheetsToTimesheetWeek(timesheetWeeks, db)
	if err != nil {
		return outData, err
	}

	outData.Data.MonthStr = monthStr
	outData.Data.Year = year
	outData.Data.Dates = dates
	outData.Data.TimesheetRows = timesheetRows
	outData.Data.WeekStartDate = weekStartDate

	return outData, nil
}

func GetAvailableJobs(employeeId, weekStartDate string, db *sql.DB) ([]domain.Job, error) {
	//TODO: validate weekStartDate

	outJobs := []domain.Job{}
	isAvailable := true
	jobs, err := infrastructure.GetJobs(db)
	if err != nil {
		return nil, err
	}

	currTimesheetWeek, err := infrastructure.GetTimesheetWeekByEmployeeWeekStart(employeeId, weekStartDate, db)
	if err == sql.ErrNoRows {
		return jobs, nil
	} else if err != nil {
		return nil, err
	}

	for i := range jobs {
		for j := range currTimesheetWeek {
			if jobs[i].ID == currTimesheetWeek[j].JobId {
				isAvailable = false
			}
		}
		if isAvailable && !jobs[i].IsComplete {
			outJobs = append(outJobs, jobs[i])
		}
		isAvailable = true
	}

	return outJobs, nil
}
