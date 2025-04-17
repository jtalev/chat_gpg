package application

import (
	"database/sql"
	"log"

	models "github.com/jtalev/chat_gpg/domain/models"
	infrastructure "github.com/jtalev/chat_gpg/infrastructure/repository"
)

type Role struct {
	db   *sql.DB
	Role models.Role
}

func (r *Role) GetEmployeeRole() (models.Role, error) {
	role, err := infrastructure.GetEmployeeRole(r.Role.EmployeeId, r.db)
	if err != nil {
		return role, err
	}
	return role, nil
}

func (r *Role) PostEmployeeRole() (models.Role, error) {
	role, err := infrastructure.PostEmployeeRole(r.Role, r.db)
	if err != nil {
		return role, err
	}
	return role, nil
}

func (r *Role) PutEmployeeRole() (models.Role, error) {
	log.Println(r.Role)
	role, err := infrastructure.PutEmployeeRole(r.Role, r.db)
	if err != nil {
		return role, err
	}
	return role, nil
}

func (r *Role) DeleteEmployeeRole() (bool, error) {
	isdeleted, err := infrastructure.DeleteEmployeeRole(r.Role.UUID, r.db)
	if err != nil {
		return false, err
	}
	return isdeleted, nil
}
