package repository

import (
	"database/sql"
	"log"

	"github.com/jtalev/chat_gpg/models"
)

func GetEmployeeByEmployeeId(employeeId string, db *sql.DB) (models.Employee, error) {
	employee := models.Employee{}
	q := `
	select * from employee where employee_id = ?;
	`
	rows, err := db.Query(q, employeeId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&employee.ID, &employee.EmployeeId, &employee.FirstName, &employee.LastName,
			&employee.Email, &employee.PhoneNumber, &employee.IsAdmin, &employee.CreatedAt, &employee.UpdatedAt)
		if err != nil {
			return employee, err
		}
	} else {
		return employee, sql.ErrNoRows
	}
	return employee, nil
}
