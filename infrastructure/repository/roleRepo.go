package infrastructure

import (
	"database/sql"
	"log"

	models "github.com/jtalev/chat_gpg/domain/models"
)

func GetEmployeeRole(employeeid string, db *sql.DB) (models.Role, error) {
	q := `
	select * from employee_role where employee_id = ?
	`

	rows, err := db.Query(q, employeeid)
	if err != nil {
		return models.Role{}, err
	}
	defer rows.Close()

	role := models.Role{}
	if rows.Next() {
		err := rows.Scan(
			&role.UUID,
			&role.EmployeeId,
			&role.Role,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return role, err
		}
	}
	return role, nil
}

func PostEmployeeRole(role models.Role, db *sql.DB) (models.Role, error) {
	q := `
	INSERT INTO employee_role (uuid, employee_id, role)
	VALUES ($1, $2, $3);
	`

	_, err := db.Exec(
		q, role.UUID, role.EmployeeId, role.Role)
	if err != nil {
		return models.Role{}, err
	}

	role, err = GetEmployeeRole(role.EmployeeId, db)
	if err != nil {
		return role, err
	}

	return role, nil
}

func PutEmployeeRole(role models.Role, db *sql.DB) (models.Role, error) {
	log.Println(role)
	q := `
	update employee_role
	set role = $1, updated_at = CURRENT_TIMESTAMP
	where employee_id = $2
	`

	_, err := db.Exec(q, role.Role, role.EmployeeId)
	if err != nil {
		return models.Role{}, err
	}
	role, err = GetEmployeeRole(role.EmployeeId, db)
	if err != nil {
		return role, err
	}
	log.Println(role)
	return role, nil
}

func DeleteEmployeeRole(uuid string, db *sql.DB) (bool, error) {
	q := `
	delete from employee_role where uuid = ?;
	`
	_, err := db.Exec(q, uuid)
	if err != nil {
		return false, err
	}

	return true, nil
}
