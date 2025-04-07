package infrastructure

import (
	"database/sql"
	"errors"
	"log"

	"github.com/jtalev/chat_gpg/domain/models"
)

func GetTimesheetWeeks(db *sql.DB) ([]domain.TimesheetWeek, error) {
	q := `
	SELECT * FROM timesheet_week;
	`
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outTimesheetWeeks := []domain.TimesheetWeek{}
	for rows.Next() {
		var timesheetWeek domain.TimesheetWeek
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

func GetTimesheetWeekById(id int, db *sql.DB) (domain.TimesheetWeek, error) {
	q := `
	SELECT * FROM timesheet_week WHERE timesheet_week_id = ?;
	`
	rows, err := db.Query(q, id)
	if err != nil {
		return domain.TimesheetWeek{}, err
	}
	defer rows.Close()

	var outTimesheetWeek domain.TimesheetWeek
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
			return domain.TimesheetWeek{}, err
		}
	} else {
		return domain.TimesheetWeek{}, sql.ErrNoRows
	}

	return outTimesheetWeek, nil
}

func GetTimesheetWeekByEmployeeWeekStart(employeeId, weekStart string, db *sql.DB) ([]domain.TimesheetWeek, error) {
	q := `
	SELECT * FROM timesheet_week where week_start_date = ? and employee_id = ?;
	`
	rows, err := db.Query(q, weekStart, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outTimesheetWeeks := []domain.TimesheetWeek{}
	for rows.Next() {
		var timesheetWeek domain.TimesheetWeek
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

func GetTimesheetWeekByWeekStart(weekStart string, db *sql.DB) ([]domain.TimesheetWeek, error) {
	q := `
	SELECT * FROM timesheet_week where week_start_date = ?;
	`
	rows, err := db.Query(q, weekStart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outTimesheetWeeks := []domain.TimesheetWeek{}
	for rows.Next() {
		var timesheetWeek domain.TimesheetWeek
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

func GetTimesheetWeekByEmployee(employeeId string, db *sql.DB) ([]domain.TimesheetWeek, error) {
	q := `
	SELECT * FROM timesheet_week where employee_id = ?;
	`
	rows, err := db.Query(q, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outTimesheetWeeks := []domain.TimesheetWeek{}
	for rows.Next() {
		var timesheetWeek domain.TimesheetWeek
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

func PostTimesheetWeek(inTimesheetWeek domain.TimesheetWeek, db *sql.DB) (domain.TimesheetWeek, error) {
	q := `
	INSERT INTO timesheet_week (employee_id, job_id, week_start_date) 
	VALUES ($1, $2, $3) 
	RETURNING timesheet_week_id, employee_id, job_id, week_start_date, created_at, modified_at;
	`

	var outTimesheetWeek domain.TimesheetWeek
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
		return domain.TimesheetWeek{}, err
	}

	return outTimesheetWeek, nil
}

func PutTimesheetWeek(inTimesheetWeek domain.TimesheetWeek, db *sql.DB) (domain.TimesheetWeek, error) {
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
		return domain.TimesheetWeek{}, err
	}
	outTimesheetWeek, err := GetTimesheetWeekById(inTimesheetWeek.TimesheetWeekId, db)
	if err != nil {
		log.Println(err)
		return domain.TimesheetWeek{}, err
	}
	return outTimesheetWeek, nil
}

func DeleteTimesheetWeek(id int, db *sql.DB) (domain.TimesheetWeek, error) {
	q := `
	delete from timesheet_week where timesheet_week_id = ?;
	`
	result, err := db.Exec(q, id)
	if err != nil {
		return domain.TimesheetWeek{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.TimesheetWeek{}, err
	}

	if rowsAffected == 0 {
		return domain.TimesheetWeek{}, errors.New("no timesheet found with the provided id")
	}

	outTimesheetWeek, err := GetTimesheetWeekById(id, db)
	if err == nil {
		return outTimesheetWeek, errors.New("TimesheetWeek still exists")
	}

	return outTimesheetWeek, nil
}
