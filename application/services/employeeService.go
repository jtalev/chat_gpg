package application

import (
	"database/sql"
	"strconv"

	"github.com/jtalev/chat_gpg/domain/models"
	"github.com/jtalev/chat_gpg/infrastructure/repository"
)

func GetEmployees(db *sql.DB) ([]domain.Employee, error) {
	outEmployees, err := infrastructure.GetEmployees(db)
	if err != nil {
		return nil, err
	}
	return outEmployees, nil
}

func GetEmployeeByEmployeeId(employeeId string, db *sql.DB) (domain.Employee, error) {
	outEmployee, err := infrastructure.GetEmployeeByEmployeeId(employeeId, db)
	if err != nil {
		return domain.Employee{}, err
	}
	return outEmployee, nil
}

func GetEmployeeById(id int, db *sql.DB) (domain.Employee, error) {
	outEmployee, err := infrastructure.GetEmployeeById(id, db)
	if err != nil {
		return domain.Employee{}, err
	}
	return outEmployee, nil
}

func DeleteAndReturnEmployees(idStr string, db *sql.DB) ([]domain.Employee, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	isDeleted, err := infrastructure.DeleteEmployee(id, db)
	if !isDeleted {
		return nil, err
	}
	employees, err := GetEmployees(db)
	if err != nil {
		return nil, err
	}
	return employees, nil
}
