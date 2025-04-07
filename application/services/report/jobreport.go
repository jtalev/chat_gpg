package report

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	domain "github.com/jtalev/chat_gpg/domain/models"
	repo "github.com/jtalev/chat_gpg/infrastructure/repository"
)

type Job struct {
	ID   int
	Name string
}

type InitDataJobReport struct {
	Jobs []Job
}

func getJobs(db *sql.DB) ([]domain.Job, error) {
	jobs, err := repo.GetJobs(db)
	if err != nil {
		log.Printf("error getting jobs: %v", err)
		return nil, err
	}
	return jobs, nil
}

func convertToJobSelect(jobs []domain.Job) []Job {
	len := len(jobs)
	out := make([]Job, len)
	for i, job := range jobs {
		out[i].ID = job.ID
		out[i].Name = fmt.Sprintf("%s", job.Name)
	}
	return out
}

func InitJobReportData(db *sql.DB) (InitDataJobReport, error) {
	var outData InitDataJobReport

	jobs, err := getJobs(db)
	if err != nil {
		return outData, err
	}

	outData.Jobs = append(outData.Jobs, convertToJobSelect(jobs)...)

	return outData, nil
}

type EmployeeTimesheetWeek struct {
	Timesheets *[]domain.Timesheet
	Hrs        string
}

type EmployeeJobReport struct {
	Name       string
	EmployeeId string
	*EmployeeTimesheetWeek
}

type JobReport struct {
	*Job
	Hrs             string
	EmployeeReports []EmployeeJobReport
}

func jobSelectMap(id int, db *sql.DB) (*Job, error) {
	jobSelect := Job{
		ID: id,
	}
	job, err := repo.GetJobById(id, db)
	if err != nil {
		log.Printf("error retrieving job from db: %v", err)
		return &jobSelect, err
	}
	jobSelect.Name = job.Name
	return &jobSelect, nil
}

func startDate(date time.Time) time.Time {
	now := time.Now()
	switch now.Weekday() {
	case time.Wednesday:
		return now
	case time.Thursday:
		now = now.AddDate(0, 0, -1)
		return now
	case time.Friday:
		now = now.AddDate(0, 0, -2)
		return now
	case time.Saturday:
		now = now.AddDate(0, 0, -3)
		return now
	case time.Sunday:
		now = now.AddDate(0, 0, -4)
		return now
	case time.Monday:
		now = now.AddDate(0, 0, -5)
		return now
	case time.Tuesday:
		now = now.AddDate(0, 0, -6)
		return now
	}
	return now
}

func startDateStr(date time.Time) string {
	now := time.Now()
	switch now.Weekday() {
	case time.Wednesday:
		return fmt.Sprintf("%v-%v-%v", now.Year(), int(now.Month()), now.Day())
	case time.Thursday:
		now = now.AddDate(0, 0, -1)
		return fmt.Sprintf("%v-%v-%v", now.Year(), int(now.Month()), now.Day())
	case time.Friday:
		now = now.AddDate(0, 0, -2)
		return fmt.Sprintf("%v-%v-%v", now.Year(), int(now.Month()), now.Day())
	case time.Saturday:
		now = now.AddDate(0, 0, -3)
		return fmt.Sprintf("%v-%v-%v", now.Year(), int(now.Month()), now.Day())
	case time.Sunday:
		now = now.AddDate(0, 0, -4)
		return fmt.Sprintf("%v-%v-%v", now.Year(), int(now.Month()), now.Day())
	case time.Monday:
		now = now.AddDate(0, 0, -5)
		return fmt.Sprintf("%v-%v-%v", now.Year(), int(now.Month()), now.Day())
	case time.Tuesday:
		now = now.AddDate(0, 0, -6)
		return fmt.Sprintf("%v-%v-%v", now.Year(), int(now.Month()), now.Day())
	}
	return ""
}

func getTimesheetWeeks(db *sql.DB) ([]domain.TimesheetWeek, error) {
	timesheetWeeks, err := repo.GetTimesheetWeeks(db)
	if err != nil {
		log.Print("error getting timesheet week at GetTimesheetWeekByWeekStart: %v", err)
		return nil, err
	}

	return timesheetWeeks, nil
}

func dateStrToDate(date string) time.Time {
	slc := strings.Split(date, "-")
	year, month, day := slc[0], slc[1], slc[2]
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		panic("year is not a number")
	}
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		panic("month is not a number")
	}
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		panic("day is not a number")
	}
	return time.Date(yearInt, time.Month(monthInt), dayInt, 0, 0, 0, 0, time.Local)
}

func filterTimesheetWeeksByDate(weeks int, timesheetWeeks []domain.TimesheetWeek) []domain.TimesheetWeek {
	out := []domain.TimesheetWeek{}

	if weeks == 0 {
		log.Println("no duration for timesheet weeks submitted")
		return nil
	}

	date := startDate(time.Now())
	date = date.AddDate(0, 0, -7*(weeks-1))

	for _, timesheetWeek := range timesheetWeeks {
		weekStartDate := dateStrToDate(timesheetWeek.WeekStartDate)
		if weekStartDate.After(date.AddDate(0, 0, -1)) {
			out = append(out, timesheetWeek)
		}
	}

	return out
}

func filterTimesheetWeeksByJob(jobId int, timesheetWeeks []domain.TimesheetWeek) []domain.TimesheetWeek {
	var filtered = make([]domain.TimesheetWeek, 0, len(timesheetWeeks))

	for _, week := range timesheetWeeks {
		if week.JobId == jobId {
			filtered = append(filtered, week)
		}
	}

	return filtered
}

func iterateTimesheetWeekFields(timesheetWeek domain.TimesheetWeek) []int {
	val := reflect.ValueOf(timesheetWeek)
	out := make([]int, 7)

	if val.Kind() != reflect.Struct {
		log.Println("interface s is not a struct")
		return out
	}

	j := 0
	for i := 3; i < val.NumField()-4; i++ {
		out[j] = int(val.Field(i).Int())
		j++
	}
	return out
}

func calcEmployeeWeekHrs(timesheets []domain.Timesheet) string {
	hrs, mins := 0, 0
	for _, timesheet := range timesheets {
		hrs += timesheet.Hours
		mins += timesheet.Minutes
	}
	hrs += mins / 60
	mins = mins % 60
	return fmt.Sprintf("%v:%v", hrs, mins)
}

func generateEmployeeReports(timesheetWeeks []domain.TimesheetWeek, db *sql.DB) (*[]EmployeeJobReport, error) {
	var employeeJobReports []EmployeeJobReport
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, week := range timesheetWeeks {
		wg.Add(1)

		go func(week domain.TimesheetWeek) {
			defer wg.Done()

			timesheetIds := iterateTimesheetWeekFields(week)
			timesheets := make([]domain.Timesheet, 7)
			name := ""
			hrs := ""

			for i, id := range timesheetIds {
				employee, err := repo.GetEmployeeByEmployeeId(week.EmployeeId, db)
				name = fmt.Sprintf("%s %s", employee.FirstName, employee.LastName)
				if err != nil {
					log.Printf("error getting employee: %v", err)
					return
				}
				timesheet, err := repo.GetTimesheetById(id, db)
				if err != nil {
					log.Printf("error getting timesheet: %v", err)
					return
				}
				timesheets[i] = timesheet
			}
			hrs = calcEmployeeWeekHrs(timesheets)

			employeeTimesheetWeek := EmployeeTimesheetWeek{
				Timesheets: &timesheets,
				Hrs:        hrs,
			}
			employeeReport := EmployeeJobReport{
				Name:                  name,
				EmployeeId:            week.EmployeeId,
				EmployeeTimesheetWeek: &employeeTimesheetWeek,
			}

			mu.Lock()
			employeeJobReports = append(employeeJobReports, employeeReport)
			mu.Unlock()
		}(week)
	}

	wg.Wait()

	return &employeeJobReports, nil
}

// func determineWeekStart(weekStartStr, arrow string) string {
// 	if weekStartStr == "" {
// 		return startDate(time.Now())
// 	}

// 	slc := strings.Split(weekStartStr, "-")
// 	if len(slc) != 3 {
// 		panic(`date string arr should always be len = 3 when split at "-"`)
// 	}
// 	yearstr, monthstr, daystr := slc[0], slc[1], slc[2]
// 	year, err := strconv.Atoi(yearstr)
// 	month, err := strconv.Atoi(monthstr)
// 	day, err := strconv.Atoi(daystr)
// 	if err != nil {
// 		panic("bad request, cannot convert date str to int")
// 	}

// 	weekStartDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
// 	if arrow == "left" {
// 		weekStartDate = weekStartDate.AddDate(0, 0, -7)
// 	} else if arrow == "right" {
// 		weekStartDate = weekStartDate.AddDate(0, 0, 7)
// 	} else {
// 		panic(`bad request, arrow should be "left" or "right"`)
// 	}
// 	weekStartStr = fmt.Sprintf("%v-%v-%v", weekStartDate.Year(), int(weekStartDate.Month()), weekStartDate.Day())
// 	log.Println(weekStartStr)
// 	return weekStartStr
// }

func GetJobReportData(jobId int, weeks int, db *sql.DB) (JobReport, error) {
	var jobReport JobReport

	jobSelect, err := jobSelectMap(jobId, db)
	if err != nil {
		return jobReport, err
	}

	timesheetWeeks, err := getTimesheetWeeks(db)
	if err != nil {
		log.Printf("error getting timesheet week: %v", err)
		return jobReport, err
	}

	filtered := filterTimesheetWeeksByJob(jobId, timesheetWeeks)
	filtered = filterTimesheetWeeksByDate(weeks, filtered)

	employeeJobReports, err := generateEmployeeReports(filtered, db)

	jobReport.Job = jobSelect
	jobReport.EmployeeReports = *employeeJobReports
	return jobReport, nil
}
