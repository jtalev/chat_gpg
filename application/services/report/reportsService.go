package report

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	application "github.com/jtalev/chat_gpg/application/services"
	domain "github.com/jtalev/chat_gpg/domain/models"
	repo "github.com/jtalev/chat_gpg/infrastructure/repository"
)

type TimesheetReportData struct {
	Employees     []domain.Employee
	WeekStartDate string
}

func InitialTimesheetReportData(db *sql.DB) (TimesheetReportData, error) {
	outData := TimesheetReportData{}
	employees, err := application.GetEmployees(db)
	if err != nil {
		return TimesheetReportData{}, err
	}
	weekStartDate := application.WeekStartDate()
	weekStartDateStr := fmt.Sprintf("%v-%v-%v", weekStartDate.Year(), int(weekStartDate.Month()), weekStartDate.Day())

	sort.Slice(employees, func(i, j int) bool {
		return employees[i].FirstName < employees[j].FirstName
	})
	outData.Employees = employees
	outData.WeekStartDate = weekStartDateStr

	return outData, nil
}

type EmployeeTimesheetReportData struct {
	EmployeeId            string
	WeekStartDate         string
	WeekDates             []int
	TimesheetRows         []application.TimesheetRow
	DayTotals             []string
	WeekTotal             string
	LeaveHoursPayable     string
	TotalHoursPayable     string
	RelevantLeaveRequests []domain.LeaveRequest
}

func calcTimesheetRowTotal(inTimesheetRows []application.TimesheetRow) []string {
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

func formatTimesheetMatrices(inTimesheetRows []application.TimesheetRow) (outHoursMatrix, outMinutesMatrix [][]int) {
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
		return []string{"0:0", "0:0", "0:0", "0:0", "0:0", "0:0", "0:0"}
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

func formatDayTotals(inTimesheetRows []application.TimesheetRow) []string {
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

func GetCurrentWeekDates(weekStartDate string) ([]int, error) {
	outData := make([]int, 7)
	dateArr := strings.Split(weekStartDate, "-")
	dateArrInt := make([]int, 3)
	for i := range dateArr {
		date, err := strconv.Atoi(dateArr[i])
		if err != nil {
			return nil, err
		}
		dateArrInt[i] = date
	}
	currentDate := time.Date(dateArrInt[0], time.Month(dateArrInt[1]), dateArrInt[2], 0, 0, 0, 0, time.Local)
	for i := range outData {
		outData[i] = currentDate.Day()
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return outData, nil
}

func GetEmployeeTimesheetReport(id, weekStartDate string, db *sql.DB) (EmployeeTimesheetReportData, error) {
	outData := EmployeeTimesheetReportData{
		WeekStartDate: weekStartDate,
	}

	idStr, err := strconv.Atoi(id)
	if err != nil {
		return outData, err
	}

	employee, err := application.GetEmployeeById(idStr, db)
	if err != nil {
		return outData, err
	}

	weekDates, err := GetCurrentWeekDates(weekStartDate)
	if err != nil {
		return outData, err
	}

	timesheetWeeks, err := application.GetTimesheetWeekByEmployeeWeekStart(employee.EmployeeId, weekStartDate, db)
	if err != nil {
		return outData, err
	}
	timesheetRows, err := application.MapTimesheetsToTimesheetWeek(timesheetWeeks, db)
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

	leaveHrsPayable, err := ProcessLeavePayable(weekStartDate, employee.EmployeeId, db)
	if err != nil {
		return outData, err
	}

	employeeLeaveRequests, err := GetLeaveByEmployeeId(employee.EmployeeId, db)
	if err != nil {
		return outData, err
	}
	relevantLeaveRequests, err := FilterEmployeeLeaveByWeek(employeeLeaveRequests, weekStartDate)
	if err != nil {
		return outData, err
	}

	outData.EmployeeId = employee.EmployeeId
	outData.WeekDates = weekDates
	outData.TimesheetRows = timesheetRows
	outData.DayTotals = dayTotals
	outData.WeekTotal = weekTotal
	outData.LeaveHoursPayable = leaveHrsPayable
	outData.RelevantLeaveRequests = relevantLeaveRequests

	totalHoursPayable, err := calcTotalHrsPayable(outData.WeekTotal, outData.LeaveHoursPayable)
	if err != nil {
		return outData, err
	}
	outData.TotalHoursPayable = totalHoursPayable

	return outData, nil
}

func GetPrevEmployeeTimesheetReport(id, weekStartDate string, db *sql.DB) (EmployeeTimesheetReportData, error) {
	weekStartDateArr := strings.Split(weekStartDate, "-")
	weekStartDateArrInt := make([]int, 3)
	for i := range weekStartDateArr {
		val, err := strconv.Atoi(weekStartDateArr[i])
		if err != nil {
			return EmployeeTimesheetReportData{}, err
		}
		weekStartDateArrInt[i] = val
	}
	date := time.Date(weekStartDateArrInt[0], time.Month(weekStartDateArrInt[1]), weekStartDateArrInt[2], 0, 0, 0, 0, time.Local)
	date = date.AddDate(0, 0, -7)
	weekStartDate = fmt.Sprintf("%v-%v-%v", date.Year(), int(date.Month()), date.Day())
	outData := EmployeeTimesheetReportData{
		WeekStartDate: weekStartDate,
	}

	weekDates, err := GetCurrentWeekDates(weekStartDate)
	if err != nil {
		return outData, err
	}

	timesheetWeeks, err := application.GetTimesheetWeekByEmployeeWeekStart(id, outData.WeekStartDate, db)
	if err != nil {
		return outData, err
	}
	timesheetRows, err := application.MapTimesheetsToTimesheetWeek(timesheetWeeks, db)
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

	leaveHrsPayable, err := ProcessLeavePayable(weekStartDate, id, db)
	if err != nil {
		return outData, err
	}

	employeeLeaveRequests, err := GetLeaveByEmployeeId(id, db)
	if err != nil {
		return outData, err
	}
	relevantLeaveRequests, err := FilterEmployeeLeaveByWeek(employeeLeaveRequests, weekStartDate)
	if err != nil {
		return outData, err
	}

	outData.EmployeeId = id
	outData.WeekDates = weekDates
	outData.TimesheetRows = timesheetRows
	outData.DayTotals = dayTotals
	outData.WeekTotal = weekTotal
	outData.LeaveHoursPayable = leaveHrsPayable
	outData.RelevantLeaveRequests = relevantLeaveRequests

	totalHoursPayable, err := calcTotalHrsPayable(outData.WeekTotal, outData.LeaveHoursPayable)
	if err != nil {
		return outData, err
	}
	outData.TotalHoursPayable = totalHoursPayable

	return outData, nil
}

func GetNextEmployeeTimesheetReport(id, weekStartDate string, db *sql.DB) (EmployeeTimesheetReportData, error) {
	weekStartDateArr := strings.Split(weekStartDate, "-")
	weekStartDateArrInt := make([]int, 3)
	for i := range weekStartDateArr {
		val, err := strconv.Atoi(weekStartDateArr[i])
		if err != nil {
			return EmployeeTimesheetReportData{}, err
		}
		weekStartDateArrInt[i] = val
	}
	date := time.Date(weekStartDateArrInt[0], time.Month(weekStartDateArrInt[1]), weekStartDateArrInt[2], 0, 0, 0, 0, time.Local)
	date = date.AddDate(0, 0, 7)
	weekStartDate = fmt.Sprintf("%v-%v-%v", date.Year(), int(date.Month()), date.Day())
	outData := EmployeeTimesheetReportData{
		WeekStartDate: weekStartDate,
	}

	weekDates, err := GetCurrentWeekDates(weekStartDate)
	if err != nil {
		return outData, err
	}

	timesheetWeeks, err := application.GetTimesheetWeekByEmployeeWeekStart(id, outData.WeekStartDate, db)
	if err != nil {
		return outData, err
	}
	timesheetRows, err := application.MapTimesheetsToTimesheetWeek(timesheetWeeks, db)
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

	leaveHrsPayable, err := ProcessLeavePayable(weekStartDate, id, db)
	if err != nil {
		return outData, err
	}

	employeeLeaveRequests, err := GetLeaveByEmployeeId(id, db)
	if err != nil {
		return outData, err
	}
	relevantLeaveRequests, err := FilterEmployeeLeaveByWeek(employeeLeaveRequests, weekStartDate)
	if err != nil {
		return outData, err
	}

	outData.EmployeeId = id
	outData.WeekDates = weekDates
	outData.TimesheetRows = timesheetRows
	outData.DayTotals = dayTotals
	outData.WeekTotal = weekTotal
	outData.LeaveHoursPayable = leaveHrsPayable
	outData.RelevantLeaveRequests = relevantLeaveRequests

	totalHoursPayable, err := calcTotalHrsPayable(outData.WeekTotal, outData.LeaveHoursPayable)
	if err != nil {
		return outData, err
	}
	outData.TotalHoursPayable = totalHoursPayable

	return outData, nil
}

func GetLeaveByEmployeeId(employeeId string, db *sql.DB) ([]domain.LeaveRequest, error) {
	leaveRequests, err := repo.GetLeaveRequestsByEmployee(employeeId, db)
	if err != nil {
		return nil, err
	}

	return leaveRequests, nil
}

func FilterEmployeeLeaveByWeek(leaveRequests []domain.LeaveRequest, weekStartDate string) ([]domain.LeaveRequest, error) {
	outLeaveRequests := []domain.LeaveRequest{}
	startDate, err := application.DateStrToDate(weekStartDate)
	if err != nil {
		return nil, err
	}
	endDate := startDate.AddDate(0, 0, 7)

	for _, lr := range leaveRequests {
		if !lr.IsApproved || lr.Type == "unpaid" {
			continue
		}

		lrStart, err := application.DateStrToDate(lr.From)
		lrEnd, err := application.DateStrToDate(lr.To)
		if err != nil {
			return nil, err
		}

		for lrStart.Before(lrEnd.AddDate(0, 0, 1)) {
			if lrStart.After(startDate.AddDate(0, 0, -1)) && lrStart.Before(endDate) {
				outLeaveRequests = append(outLeaveRequests, lr)
				break
			}
			lrStart = lrStart.AddDate(0, 0, 1)
		}

	}

	return outLeaveRequests, nil
}

func CalcLeavePayable(filteredLeaveRequests []domain.LeaveRequest, weekStartDate string) (string, error) {
	hrs := 0

	startDate, err := application.DateStrToDate(weekStartDate)
	if err != nil {
		return "", nil
	}
	endDate := startDate.AddDate(0, 0, 7)

	for _, lr := range filteredLeaveRequests {
		lrStart, err := application.DateStrToDate(lr.From)
		lrEnd, err := application.DateStrToDate(lr.To)
		if err != nil {
			return "", err
		}

		for lrStart.Before(lrEnd.AddDate(0, 0, 1)) {
			if lrStart.Weekday() == time.Saturday || lrStart.Weekday() == time.Sunday {
				lrStart = lrStart.AddDate(0, 0, 1)
				continue
			}
			if lrStart.After(startDate.AddDate(0, 0, -1)) && lrStart.Before(endDate) {
				hrs += lr.HoursPerDay
			}
			lrStart = lrStart.AddDate(0, 0, 1)
		}
	}

	return strconv.Itoa(hrs), nil
}

func ProcessLeavePayable(weekStartDate, employeeId string, db *sql.DB) (string, error) {
	employeesLeaveRequests, err := GetLeaveByEmployeeId(employeeId, db)
	if err != nil {
		return "", err
	}

	filteredLeaveRequests, err := FilterEmployeeLeaveByWeek(employeesLeaveRequests, weekStartDate)
	if err != nil {
		return "", err
	}

	leavePayableStr, err := CalcLeavePayable(filteredLeaveRequests, weekStartDate)
	if err != nil {
		return "", err
	}

	return leavePayableStr, nil
}

func calcTotalHrsPayable(weekTotal, leaveHoursPayable string) (string, error) {
	week := strings.Split(weekTotal, ":")

	weekInt := make([]int, 2)
	hrs, mins := 0, 0

	for i := range week {
		val, err := strconv.Atoi(week[i])
		if err != nil {
			return "", err
		}
		weekInt[i] = val
	}

	leaveHrs, err := strconv.Atoi(leaveHoursPayable)
	if err != nil {
		return "", nil
	}

	hrs = weekInt[0] + leaveHrs
	mins = weekInt[1]

	hrs += mins / 60
	mins = mins % 60

	return fmt.Sprintf("%v:%v", hrs, mins), nil
}
