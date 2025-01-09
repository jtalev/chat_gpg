package services

import (
	"database/sql"

	"github.com/jtalev/chat_gpg/models"
	"github.com/jtalev/chat_gpg/repository"
)

func GetEmployeeByEmployeeId(employeeId string, db *sql.DB) (models.Employee, error) {
	outEmployee, err := repository.GetEmployeeByEmployeeId(employeeId, db)
	if err != nil {
		return models.Employee{}, err
	}
	return outEmployee, nil
}
