package infrastructure

import (
	"database/sql"
	"log"

	"github.com/jtalev/chat_gpg/domain/models"
)

func GetEmployees(db *sql.DB) ([]domain.Employee, error) {
	q := `
	select * from employee;
	`
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outEmployees := []domain.Employee{}
	for rows.Next() {
		employee := domain.Employee{}
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

func GetEmployeeById(id int, db *sql.DB) (domain.Employee, error) {
	q := `
	select * from employee where id = ?;
	`
	rows, err := db.Query(q, id)
	if err != nil {
		return domain.Employee{}, err
	}
	defer rows.Close()

	outEmployee := domain.Employee{}
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
			return domain.Employee{}, err
		}
	}
	return outEmployee, nil

}

func GetEmployeeByEmployeeId(employeeId string, db *sql.DB) (domain.Employee, error) {
	employee := domain.Employee{}
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

func DeleteEmployee(id int, db *sql.DB) (bool, error) {
	q := `
	delete from employee where id = ?;
	`

	_, err := db.Exec(q, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func PostEmployee(employee domain.Employee, db *sql.DB) (domain.Employee, error) {
	q := `
	INSERT INTO employee (employee_id, first_name, last_name, email, phone_number, is_admin)
	VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err := db.Exec(q, employee.EmployeeId, employee.FirstName, employee.LastName, employee.Email, employee.PhoneNumber, employee.IsAdmin)
	if err != nil {
		return domain.Employee{}, err
	}

	outEmployee, err := GetEmployeeByEmployeeId(employee.EmployeeId, db)
	if err != nil {
		return domain.Employee{}, err
	}

	return outEmployee, nil
}
