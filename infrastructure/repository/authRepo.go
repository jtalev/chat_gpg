package infrastructure

import (
	"database/sql"
	"log"

	"github.com/jtalev/chat_gpg/domain/models"
)

func GetEmployeeAuthByUsername(username string, db *sql.DB) (domain.EmployeeAuth, error) {
	employeeAuth := domain.EmployeeAuth{}
	q := `
	select * from employee_auth where username = ?;
	`
	rows, err := db.Query(q, username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&employeeAuth.AuthId, &employeeAuth.EmployeeId, &employeeAuth.Username,
			&employeeAuth.PasswordHash, &employeeAuth.CreatedAt, &employeeAuth.UpdatedAt)
		if err != nil {
			return employeeAuth, err
		}
	} else {
		return employeeAuth, sql.ErrNoRows
	}
	return employeeAuth, nil
}

func PostEmployeeAuth(employeeAuth domain.EmployeeAuth, db *sql.DB) (domain.EmployeeAuth, error) {
	q := `
	INSERT INTO employee_auth (employee_id, username, password_hash)
	VALUES ($1, $2, $3);
	`

	_, err := db.Exec(q, employeeAuth.EmployeeId, employeeAuth.Username, employeeAuth.PasswordHash)
	if err != nil {
		return domain.EmployeeAuth{}, err
	}

	outEmployee, err := GetEmployeeAuthByUsername(employeeAuth.Username, db)
	if err != nil {
		return domain.EmployeeAuth{}, err
	}

	return outEmployee, nil
}
