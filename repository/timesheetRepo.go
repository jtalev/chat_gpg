package repository

import (
	"database/sql"

	"github.com/jtalev/chat_gpg/models"
)

func GetTimesheetById(id int, db *sql.DB) (models.Timesheet, error) {
	q := `
	select * from timesheet where id = ?;
	`

	rows, err := db.Query(q, id)
	if err != nil {
		return models.Timesheet{}, err
	}
	defer rows.Close()

	var timesheet models.Timesheet
	if rows.Next() {
		err := rows.Scan(
			&timesheet.ID,
			&timesheet.EmployeeId,
			&timesheet.JobId,
			&timesheet.WeekStart,
			&timesheet.Date,
			&timesheet.Hours,
			&timesheet.Minutes,
			&timesheet.CreatedAt,
			&timesheet.UpdatedAt,
		)
		if err != nil {
			return models.Timesheet{}, err
		}
	} else {
		return models.Timesheet{}, sql.ErrNoRows
	}

	return timesheet, nil
}

func GetTimesheetsByWeekStart(employeeId, weekStart string, db *sql.DB) ([]models.Timesheet, error) {
	q := `
	select *
	from timesheet
	where employee_id = $1 and week_start = $2
	order by job_id and timesheet_date;
	`

	rows, err := db.Query(q, employeeId, weekStart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []models.Timesheet
	for rows.Next() {
		var timesheet models.Timesheet
		err := rows.Scan(
			&timesheet.ID,
			&timesheet.EmployeeId,
			&timesheet.JobId,
			&timesheet.WeekStart,
			&timesheet.Date,
			&timesheet.Hours,
			&timesheet.Minutes,
			&timesheet.CreatedAt,
			&timesheet.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, timesheet)
	}
	return data, nil
}

func PostTimesheet(timesheet models.Timesheet, db *sql.DB) (models.Timesheet, error) {
	q := `
	INSERT INTO timesheet (employee_id, job_id, week_start, timesheet_date, hours, minutes)
	VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err := db.Exec(
		q,
		timesheet.EmployeeId,
		timesheet.JobId,
		timesheet.WeekStart,
		timesheet.Date,
		timesheet.Hours,
		timesheet.Minutes,
	)
	if err != nil {
		return models.Timesheet{}, err
	}
	return timesheet, nil
}

func PutTimesheet(timesheet models.Timesheet, db *sql.DB) (models.Timesheet, error) {
	q := `
	update timesheet
	set hours = $1, minutes = $2, updated_at = CURRENT_TIMESTAMP
	where id = $3;
	`

	_, err := db.Exec(q, timesheet.Hours, timesheet.Minutes, timesheet.ID)
	if err != nil {
		return models.Timesheet{}, err
	}
	ts, err := GetTimesheetById(timesheet.ID, db)
	if err != nil {
		return models.Timesheet{}, err
	}
	return ts, nil
}
