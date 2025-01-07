package repository

import (
	"database/sql"
	"log"

	"github.com/jtalev/chat_gpg/models"
)

func PostTimesheetWeek(inTimesheetWeek models.TimesheetWeek, db *sql.DB) (models.TimesheetWeek, error) {
	q := `
	INSERT INTO timesheet_week (employee_id, job_id, week_start_date) 
	VALUES ($1, $2, $3) 
	RETURNING timesheet_week_id, employee_id, job_id, week_start_date, created_at, modified_at;
	`

	var outTimesheetWeek models.TimesheetWeek
	err := db.QueryRow(
		q,
		inTimesheetWeek.EmployeeId,
		inTimesheetWeek.JobId,
		inTimesheetWeek.WeekStartDate,
	).Scan(
		&outTimesheetWeek.TimesheetWeekId,
		&outTimesheetWeek.EmployeeId,
		&outTimesheetWeek.JobId,
		&outTimesheetWeek.WeekStartDate,
		&outTimesheetWeek.CreatedAt,
		&outTimesheetWeek.ModifiedAt,
	)
	if err != nil {
		return models.TimesheetWeek{}, err
	}
	log.Println(outTimesheetWeek)
	return outTimesheetWeek, nil
}
