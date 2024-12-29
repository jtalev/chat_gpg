package repository

import (
	"database/sql"

	"github.com/jtalev/chat_gpg/models"
)

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
