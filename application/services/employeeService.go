package application

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/jtalev/chat_gpg/domain/models"
	auth "github.com/jtalev/chat_gpg/infrastructure/auth"
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

func PostAndReturnEmployees(hxVals []string, db *sql.DB) ([]domain.Employee, error) {
	employeeId := hxVals[0]
	firstName := hxVals[1]
	lastName := hxVals[2]
	email := hxVals[3]
	phoneNumber := hxVals[4]
	isAdminStr := hxVals[5]
	var isAdmin bool
	if isAdminStr == "true" {
		isAdmin = true
	} else {
		isAdmin = false
	}
	username := hxVals[6]
	password := hxVals[7]
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	employeeAuth := domain.EmployeeAuth{
		EmployeeId:   employeeId,
		Username:     username,
		PasswordHash: passwordHash,
	}

	employeeAuth, err = infrastructure.PostEmployeeAuth(employeeAuth, db)
	if err != nil {
		return nil, err
	}
	log.Println(employeeAuth)

	employee := domain.Employee{
		EmployeeId:  employeeId,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		PhoneNumber: phoneNumber,
		IsAdmin:     isAdmin,
	}

	employee, err = infrastructure.PostEmployee(employee, db)
	if err != nil {
		return nil, err
	}

	employees, err := GetEmployees(db)
	if err != nil {
		return nil, err
	}

	return employees, nil
}
