package repository

import (
	"database/sql"
	"log"

	"github.com/jtalev/chat_gpg/models"
)

func GetEmployees(db *sql.DB) ([]models.Employee, error) {
	q := `
	select * from employee;
	`
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outEmployees := []models.Employee{}
	for rows.Next() {
		employee := models.Employee{}
		if err := rows.Scan(
			&employee.ID,
			&employee.EmployeeId,
			&employee.FirstName,
			&employee.LastName,
			&employee.Email,
			&employee.PhoneNumber,
			&employee.IsAdmin,
			&employee.CreatedAt,
			&employee.UpdatedAt,
		); err != nil {
			return nil, err
		}
		outEmployees = append(outEmployees, employee)
	}
	return outEmployees, nil
}

func GetEmployeeById(id int, db *sql.DB) (models.Employee, error) {
	q := `
	select * from employee where id = ?;
	`
	rows, err := db.Query(q, id)
	if err != nil {
		return models.Employee{}, err
	}
	defer rows.Close()

	outEmployee := models.Employee{}
	if rows.Next() {
		err := rows.Scan(
			&outEmployee.ID,
			&outEmployee.EmployeeId,
			&outEmployee.FirstName,
			&outEmployee.LastName,
			&outEmployee.Email,
			&outEmployee.PhoneNumber,
			&outEmployee.IsAdmin,
			&outEmployee.CreatedAt,
			&outEmployee.UpdatedAt,
		)
		if err != nil {
			return models.Employee{}, err
		}
	}
	return outEmployee, nil

}

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
