package application

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	domain "github.com/jtalev/chat_gpg/domain/models"
)

type TimesheetReportData struct {
	Employees     []domain.Employee
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

	sort.Slice(employees, func(i, j int) bool {
		return employees[i].FirstName < employees[j].FirstName
	})
	outData.Employees = employees
	outData.WeekStartDate = weekStartDateStr

	return outData, nil
}

type EmployeeTimesheetReportData struct {
	EmployeeId    string
	WeekStartDate string
	WeekDates     []int
	TimesheetRows []TimesheetRow
	DayTotals     []string
	WeekTotal     string
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

	return outTotals
}

func formatTimesheetMatrices(inTimesheetRows []TimesheetRow) (outHoursMatrix, outMinutesMatrix [][]int) {
	row, col := len(inTimesheetRows), 7

	outHoursMatrix, outMinutesMatrix = make([][]int, row), make([][]int, row)
	for i := range row {
		outHoursMatrix[i] = make([]int, col)
		outMinutesMatrix[i] = make([]int, col)
	}

	for i := range inTimesheetRows {
		for j, timesheet := range inTimesheetRows[i].Timesheets {
			outHoursMatrix[i][j] = timesheet.Hours
			outMinutesMatrix[i][j] = timesheet.Minutes
		}
	}

	return outHoursMatrix, outMinutesMatrix
}

func calcDayTotals(inHoursMatrix, inMinutesMatrix [][]int) (dayTotals []string) {
	if len(inHoursMatrix) == 0 || len(inHoursMatrix[0]) == 0 {
		return []string{} // Return an empty slice if there's no data
	}

	hourTotals := make([]int, 7)
	minuteTotals := make([]int, 7)
	dayTotals = make([]string, 7)
	row, col := len(inHoursMatrix), len(inHoursMatrix[0])

	for j := range col {
		for i := range row {
			hourTotals[j] += inHoursMatrix[i][j]
			minuteTotals[j] += inMinutesMatrix[i][j]
		}
	}

	for i := range col {
		hourTotals[i] += minuteTotals[i] / 60
		minuteTotals[i] = minuteTotals[i] % 60
		dayTotalStr := fmt.Sprintf("%v:%v", hourTotals[i], minuteTotals[i])
		dayTotals[i] = dayTotalStr
	}

	return dayTotals
}

func formatDayTotals(inTimesheetRows []TimesheetRow) []string {
	hoursMatrix, minutesMatrix := formatTimesheetMatrices(inTimesheetRows)
	return calcDayTotals(hoursMatrix, minutesMatrix)
}

func calcWeekTotal(rowTotals, dayTotals []string) (string, error) {
	fromRowTotalsHours, fromRowTotalsMins := 0, 0
	fromDayTotalsHours, fromDayTotalsMins := 0, 0

	for i := range rowTotals {
		arr := strings.Split(rowTotals[i], ":")
		hours, mins := arr[0], arr[1]
		hourInt, err := strconv.Atoi(hours)
		if err != nil {
			return "", err
		}
		minsInt, err := strconv.Atoi(mins)
		if err != nil {
			return "", err
		}
		fromRowTotalsHours += hourInt
		fromRowTotalsMins += minsInt
	}

	for i := range dayTotals {
		arr := strings.Split(dayTotals[i], ":")
		hours, mins := arr[0], arr[1]
		hourInt, err := strconv.Atoi(hours)
		if err != nil {
			return "", err
		}
		minsInt, err := strconv.Atoi(mins)
		if err != nil {
			return "", err
		}
		fromDayTotalsHours += hourInt
		fromDayTotalsMins += minsInt
	}

	fromRowTotalsHours += fromRowTotalsMins / 60
	fromRowTotalsMins = fromRowTotalsMins % 60
	fromDayTotalsHours += fromDayTotalsMins / 60
	fromDayTotalsMins = fromDayTotalsMins % 60

	fromRowTotal := fmt.Sprintf("%v:%v", fromRowTotalsHours, fromRowTotalsMins)
	fromDayTotal := fmt.Sprintf("%v:%v", fromDayTotalsHours, fromDayTotalsMins)

	if fromRowTotal != fromDayTotal {
		return "", errors.New("Mismatch between row and day total, bad logic")
	}

	return fromRowTotal, nil
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
	rowTotals := calcTimesheetRowTotal(timesheetRows)
	for i := range timesheetRows {
		timesheetRows[i].Total = rowTotals[i]
	}

	dayTotals := formatDayTotals(timesheetRows)
	weekTotal, err := calcWeekTotal(rowTotals, dayTotals)
	if err != nil {
		return outData, err
	}

	outData.EmployeeId = employee.EmployeeId
	outData.WeekDates = weekDates
	outData.TimesheetRows = timesheetRows
	outData.DayTotals = dayTotals
	outData.WeekTotal = weekTotal

	return outData, nil
}
