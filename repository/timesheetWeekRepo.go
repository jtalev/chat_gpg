package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/jtalev/chat_gpg/models"
)

func GetTimesheetWeekById(id int, db *sql.DB) (models.TimesheetWeek, error) {
	q := `
	SELECT * FROM timesheet_week WHERE timesheet_week_id = ?;
	`
	rows, err := db.Query(q, id)
	if err != nil {
		return models.TimesheetWeek{}, err
	}
	defer rows.Close()

	var outTimesheetWeek models.TimesheetWeek
	if rows.Next() {
		err := rows.Scan(
			&outTimesheetWeek.TimesheetWeekId,
			&outTimesheetWeek.EmployeeId,
			&outTimesheetWeek.JobId,
			&outTimesheetWeek.WedTimesheetId,
			&outTimesheetWeek.ThuTimesheetId,
			&outTimesheetWeek.FriTimesheetId,
			&outTimesheetWeek.SatTimesheetId,
			&outTimesheetWeek.SunTimesheetId,
			&outTimesheetWeek.MonTimesheetId,
			&outTimesheetWeek.TueTimesheetId,
			&outTimesheetWeek.WeekStartDate,
			&outTimesheetWeek.CreatedAt,
			&outTimesheetWeek.ModifiedAt,
		)
		if err != nil {
			return models.TimesheetWeek{}, err
		}
	} else {
		return models.TimesheetWeek{}, sql.ErrNoRows
	}

	return outTimesheetWeek, nil
}

func GetTimesheetWeekByWeekStart(weekStart string, db *sql.DB) ([]models.TimesheetWeek, error) {
	q := `
	SELECT * FROM timesheet_week where week_start_date = ?;
	`
	rows, err := db.Query(q, weekStart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outTimesheetWeeks := []models.TimesheetWeek{}
	for rows.Next() {
		var timesheetWeek models.TimesheetWeek
		if err := rows.Scan(
			&timesheetWeek.TimesheetWeekId,
			&timesheetWeek.EmployeeId,
			&timesheetWeek.JobId,
			&timesheetWeek.WedTimesheetId,
			&timesheetWeek.ThuTimesheetId,
			&timesheetWeek.FriTimesheetId,
			&timesheetWeek.SatTimesheetId,
			&timesheetWeek.SunTimesheetId,
			&timesheetWeek.MonTimesheetId,
			&timesheetWeek.TueTimesheetId,
			&timesheetWeek.WeekStartDate,
			&timesheetWeek.CreatedAt,
			&timesheetWeek.ModifiedAt,
		); err != nil {
			return nil, err
		}
		outTimesheetWeeks = append(outTimesheetWeeks, timesheetWeek)
	}
	return outTimesheetWeeks, nil
}

func GetTimesheetWeekByEmployee(employeeId string, db *sql.DB) ([]models.TimesheetWeek, error) {
	q := `
	SELECT * FROM timesheet_week where employee_id = ?;
	`
	rows, err := db.Query(q, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outTimesheetWeeks := []models.TimesheetWeek{}
	for rows.Next() {
		var timesheetWeek models.TimesheetWeek
		if err := rows.Scan(
			&timesheetWeek.TimesheetWeekId,
			&timesheetWeek.EmployeeId,
			&timesheetWeek.JobId,
			&timesheetWeek.WedTimesheetId,
			&timesheetWeek.ThuTimesheetId,
			&timesheetWeek.FriTimesheetId,
			&timesheetWeek.SatTimesheetId,
			&timesheetWeek.SunTimesheetId,
			&timesheetWeek.MonTimesheetId,
			&timesheetWeek.TueTimesheetId,
			&timesheetWeek.WeekStartDate,
			&timesheetWeek.CreatedAt,
			&timesheetWeek.ModifiedAt,
		); err != nil {
			return nil, err
		}
		outTimesheetWeeks = append(outTimesheetWeeks, timesheetWeek)
	}
	return outTimesheetWeeks, nil
}

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

	return outTimesheetWeek, nil
}

func PutTimesheetWeek(inTimesheetWeek models.TimesheetWeek, db *sql.DB) (models.TimesheetWeek, error) {
	q := `
	UPDATE timesheet_week
	SET wed_timesheet_id = $1, thu_timesheet_id = $2, fri_timesheet_id = $3,
	sat_timesheet_id = $4, sun_timesheet_id = $5, mon_timesheet_id = $6,
	tue_timesheet_id = $7, 
	modified_at = CURRENT_TIMESTAMP
	WHERE timesheet_week_id = $8;
	`

	_, err := db.Exec(
		q,
		inTimesheetWeek.WedTimesheetId,
		inTimesheetWeek.ThuTimesheetId,
		inTimesheetWeek.FriTimesheetId,
		inTimesheetWeek.SatTimesheetId,
		inTimesheetWeek.SunTimesheetId,
		inTimesheetWeek.MonTimesheetId,
		inTimesheetWeek.TueTimesheetId,
		inTimesheetWeek.TimesheetWeekId,
	)
	if err != nil {
		log.Println(err)
		return models.TimesheetWeek{}, err
	}
	outTimesheetWeek, err := GetTimesheetWeekById(inTimesheetWeek.TimesheetWeekId, db)
	if err != nil {
		log.Println(err)
		return models.TimesheetWeek{}, err
	}
	return outTimesheetWeek, nil
}

func DeleteTimesheetWeek(id int, db *sql.DB) (models.TimesheetWeek, error) {
	q := `
	delete from timesheet_week where timesheet_week_id = ?;
	`
	_, err := db.Exec(q, id)
	if err != nil {
		return models.TimesheetWeek{}, err
	}

	outTimesheetWeek, err := GetTimesheetWeekById(id, db)
	if err == nil {
		return outTimesheetWeek, errors.New("TimesheetWeek still exists")
	}

	return outTimesheetWeek, nil
}
