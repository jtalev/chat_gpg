package services

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/jtalev/chat_gpg/models"
	"github.com/jtalev/chat_gpg/repository"
)

var nilTimesheet = models.Timesheet{
	TimesheetId:     0,
	TimesheetWeekId: 0,
	TimesheetDate:   "",
	Day:             "",
	Hours:           0,
	Minutes:         0,
	CreatedAt:       "",
	ModifiedAt:      "",
}

func splitTime(inTime string) ([]string, error) {
	containsColon := strings.Contains(inTime, ":")
	if !containsColon {
		return []string{}, errors.New("Input string doesn't contain ':'")
	}
	outTimeArr := strings.Split(inTime, ":")
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

func GetAllTimesheets(db *sql.DB) ([]models.Timesheet, error) {
	outTimesheets, err := repository.GetTimesheets(db)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return outTimesheets, nil
}

func GetTimesheetById(inTimesheetId string, db *sql.DB) (models.Timesheet, error) {
	inId, err := strconv.Atoi(inTimesheetId)
	if err != nil {
		log.Println(err)
		return models.Timesheet{}, err
	}

	outTimesheet, err := repository.GetTimesheetById(inId, db)
	if err != nil {
		log.Println(err)
		return models.Timesheet{}, err
	}

	return outTimesheet, nil
}

func PostTimesheet(inTimesheet models.Timesheet, db *sql.DB) (models.Timesheet, error) {
	if inTimesheet.TimesheetId == 0 {
		inTimesheet = nilTimesheet
	}

	// TODO: validate timesheet

	outTimesheet, err := repository.PostTimesheet(inTimesheet, db)
	if err != nil {
		log.Println(err)
		return outTimesheet, err
	}

	return outTimesheet, err
}

func PutTimesheet(id int, time string, db *sql.DB) (models.Timesheet, error) {

	parsedTime, err := splitTime(time)
	if err != nil {
		return models.Timesheet{}, nil
	}

	parsedTimeInt, err := strToInt(parsedTime)
	if err != nil {
		return models.Timesheet{}, err
	}

	inTimesheet := models.Timesheet{
		TimesheetId: id,
		Hours:       parsedTimeInt[0],
		Minutes:     parsedTimeInt[1],
	}

	// TODO: validate timesheet

	outTimesheet, err := repository.PutTimesheet(inTimesheet, db)
	if err != nil {
		log.Println(err)
		return outTimesheet, err
	}

	return outTimesheet, err
}

func InitTimesheetWeek(employeeId string, jobId int, weekStartDate string, db *sql.DB) (models.TimesheetWeek, error) {
	inTimesheetWeek := models.TimesheetWeek{
		EmployeeId:    employeeId,
		JobId:         jobId,
		WeekStartDate: weekStartDate,
	}

	tempTimesheetWeek, err := repository.PostTimesheetWeek(inTimesheetWeek, db)
	if err != nil {
		return models.TimesheetWeek{}, nil
	}

	// TODO: initialize values for Day and TimesheetDate
	var inNilTimesheets = make([]models.Timesheet, 7)
	for i := 0; i < 7; i++ {
		nilTimesheet.TimesheetWeekId = tempTimesheetWeek.TimesheetWeekId
		outTimesheet, err := repository.PostTimesheet(nilTimesheet, db)
		if err != nil {
			return models.TimesheetWeek{}, err
		}
		inNilTimesheets[i] = outTimesheet
	}

	log.Println(inNilTimesheets)

	return tempTimesheetWeek, nil
}
