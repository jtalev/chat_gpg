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

func GetEmployeeAuthByEmployeeId(employeeId string, db *sql.DB) (domain.EmployeeAuth, error) {
	employeeAuth := domain.EmployeeAuth{}
	q := `
	select * from employee_auth where employee_id = ?;
	`
	rows, err := db.Query(q, employeeId)
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

func PutEmployeeAuth(employeeAuth domain.EmployeeAuth, db *sql.DB) (domain.EmployeeAuth, error) {
	q := `
	update employee_auth
	set username = $1, password_hash = $2, updated_at = CURRENT_TIMESTAMP
	where employee_id = $3
	`

	_, err := db.Exec(q, employeeAuth.Username, employeeAuth.PasswordHash, employeeAuth.EmployeeId)
	if err != nil {
		return domain.EmployeeAuth{}, err
	}
	outEmployeeAuth, err := GetEmployeeAuthByEmployeeId(employeeAuth.EmployeeId, db)
	if err != nil {
		return domain.EmployeeAuth{}, err
	}
	return outEmployeeAuth, nil
}

func DeleteEmployeeAuth(id int, db *sql.DB) (bool, error) {
	q := `
	delete from employee_auth where auth_id = ?;
	`

	_, err := db.Exec(q, id)
	if err != nil {
		return false, err
	}

	return true, nil
}
