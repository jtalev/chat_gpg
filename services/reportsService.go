package services

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/jtalev/chat_gpg/models"
)

type TimesheetReportData struct {
	Employees     []models.Employee
	WeekStartDate string
}

func InitialTimesheetReportData(db *sql.DB) (TimesheetReportData, error) {
	outData := TimesheetReportData{}
	employees, err := GetEmployees(db)
	if err != nil {
		return TimesheetReportData{}, err
	}
	weekStartDate := weekStartDate()
	weekStartDateStr := fmt.Sprintf("%v-%v-%v", weekStartDate.Year(), int(weekStartDate.Month()), weekStartDate.Day())

	outData.Employees = employees
	outData.WeekStartDate = weekStartDateStr

	return outData, nil
}

type EmployeeTimesheetReportData struct {
	EmployeeId    string
	WeekStartDate string
	WeekDates     []int
	TimesheetRows []TimesheetRow
}

func calcTimesheetRowTotal(inTimesheetRows []TimesheetRow) []string {
	outTotals := make([]string, len(inTimesheetRows))

	for i := range inTimesheetRows {
		hours, mins := 0, 0
		for _, timesheet := range inTimesheetRows[i].Timesheets {
			hours += timesheet.Hours
			mins += timesheet.Minutes
		}
		additionalHours := mins / 60
		mins = mins % 60
		hours += additionalHours
		outTotal := fmt.Sprintf("%v:%v", hours, mins)
		outTotals[i] = outTotal
	}

	log.Println(outTotals)
	return outTotals
}

func GetInitialEmployeeTimesheetReport(id, weekStartDate string, db *sql.DB) (EmployeeTimesheetReportData, error) {
	outData := EmployeeTimesheetReportData{
		WeekStartDate: weekStartDate,
	}

	idStr, err := strconv.Atoi(id)
	if err != nil {
		return outData, err
	}

	employee, err := GetEmployeeById(idStr, db)
	if err != nil {
		return outData, err
	}

	weekDates, err := currentWeekDates()
	if err != nil {
		return outData, err
	}

	timesheetWeeks, err := GetTimesheetWeekByEmployeeWeekStart(employee.EmployeeId, weekStartDate, db)
	if err != nil {
		return outData, err
	}
	timesheetRows, err := mapTimesheetsToTimesheetWeek(timesheetWeeks, db)
	if err != nil {
		return outData, err
	}
	totals := calcTimesheetRowTotal(timesheetRows)
	for i := range timesheetRows {
		timesheetRows[i].Total = totals[i]
	}

	outData.EmployeeId = employee.EmployeeId
	outData.WeekDates = weekDates
	outData.TimesheetRows = timesheetRows
	log.Println(timesheetWeeks)

	return outData, nil
}
