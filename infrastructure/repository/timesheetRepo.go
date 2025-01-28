package infrastructure

import (
	"database/sql"

	"github.com/jtalev/chat_gpg/domain/models"
)

func GetTimesheets(db *sql.DB) ([]domain.Timesheet, error) {
	q := `
	select * from timesheet;
	`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var outTimesheets []domain.Timesheet
	for rows.Next() {
		var timesheet domain.Timesheet
		err := rows.Scan(
			&timesheet.TimesheetId,
			&timesheet.TimesheetWeekId,
			&timesheet.TimesheetDate,
			&timesheet.Day,
			&timesheet.Hours,
			&timesheet.Minutes,
			&timesheet.CreatedAt,
			&timesheet.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}
		outTimesheets = append(outTimesheets, timesheet)
	}
	return outTimesheets, nil
}

func GetTimesheetById(id int, db *sql.DB) (domain.Timesheet, error) {
	q := `
	select * from timesheet where timesheet_id = ?;
	`

	rows, err := db.Query(q, id)
	if err != nil {
		return domain.Timesheet{}, err
	}
	defer rows.Close()

	var outTimesheet domain.Timesheet
	if rows.Next() {
		err := rows.Scan(
			&outTimesheet.TimesheetId,
			&outTimesheet.TimesheetWeekId,
			&outTimesheet.TimesheetDate,
			&outTimesheet.Day,
			&outTimesheet.Hours,
			&outTimesheet.Minutes,
			&outTimesheet.CreatedAt,
			&outTimesheet.ModifiedAt,
		)
		if err != nil {
			return domain.Timesheet{}, err
		}
	} else {
		return domain.Timesheet{}, sql.ErrNoRows
	}

	return outTimesheet, nil
}

func PostTimesheet(inTimesheet domain.Timesheet, db *sql.DB) (domain.Timesheet, error) {
	q := `
	INSERT INTO timesheet (timesheet_week_id, timesheet_date, day, hours, minutes)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING timesheet_id, timesheet_week_id, timesheet_date, day, hours, minutes, created_at, modified_at;
	`

	var outTimesheet domain.Timesheet
	err := db.QueryRow(
		q,
		inTimesheet.TimesheetWeekId,
		inTimesheet.TimesheetDate,
		inTimesheet.Day,
		inTimesheet.Hours,
		inTimesheet.Minutes,
	).Scan(
		&outTimesheet.TimesheetId,
		&outTimesheet.TimesheetWeekId,
		&outTimesheet.TimesheetDate,
		&outTimesheet.Day,
		&outTimesheet.Hours,
		&outTimesheet.Minutes,
		&outTimesheet.CreatedAt,
		&outTimesheet.ModifiedAt,
	)
	if err != nil {
		return domain.Timesheet{}, err
	}
	return outTimesheet, nil
}

func PutTimesheet(inTimesheet domain.Timesheet, db *sql.DB) (domain.Timesheet, error) {
	q := `
	update timesheet
	set hours = $1, minutes = $2, modified_at = CURRENT_TIMESTAMP
	where timesheet_id = $3;
	`

	_, err := db.Exec(q, inTimesheet.Hours, inTimesheet.Minutes, inTimesheet.TimesheetId)
	if err != nil {
		return domain.Timesheet{}, err
	}
	outTimesheet, err := GetTimesheetById(inTimesheet.TimesheetId, db)
	if err != nil {
		return domain.Timesheet{}, err
	}
	return outTimesheet, nil
}
