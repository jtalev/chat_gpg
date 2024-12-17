package handlers

import (
	"database/sql"
	"log"
)

type Employee struct {
	ID          int    `json:"id"`
	EmployeeId  string `json:"employee_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	IsAdmin     bool   `json:"is_admin"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func GetEmployeeByEmployeeId(employeeId int, db *sql.DB) (Employee, error) {
	employee := Employee{}
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
