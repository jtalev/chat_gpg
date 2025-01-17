package services

import (
	"database/sql"

	"github.com/jtalev/chat_gpg/models"
	"github.com/jtalev/chat_gpg/repository"
)

func GetEmployees(db *sql.DB) ([]models.Employee, error) {
	outEmployees, err := repository.GetEmployees(db)
	if err != nil {
		return nil, err
	}
	return outEmployees, nil
}

func GetEmployeeByEmployeeId(employeeId string, db *sql.DB) (models.Employee, error) {
	outEmployee, err := repository.GetEmployeeByEmployeeId(employeeId, db)
	if err != nil {
		return models.Employee{}, err
	}
	return outEmployee, nil
}

func GetEmployeeById(id int, db *sql.DB) (models.Employee, error) {
	outEmployee, err := repository.GetEmployeeById(id, db)
	if err != nil {
		return models.Employee{}, err
	}
	return outEmployee, nil
}
