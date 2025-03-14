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

type EmployeeDto struct {
	ID             string
	EmployeeId     string
	FirstName      string
	LastName       string
	Username       string
	Password       string
	Email          string
	PhoneNumber    string
	IsAdmin        string
	EmployeeIdErr  string
	FirstNameErr   string
	LastNameErr    string
	UsernameErr    string
	PasswordErr    string
	EmailErr       string
	PhoneNumberErr string
	SuccessMsg     string
}

func employeeDtoToEmployee(employeeDto EmployeeDto) (domain.Employee, error) {
	id, err := strconv.Atoi(employeeDto.ID)
	if err != nil {
		log.Println("strconv")
		return domain.Employee{}, err
	}
	isAdmin := false
	if employeeDto.IsAdmin == "true" {
		isAdmin = true
	}
	employee := domain.Employee{
		ID:          id,
		EmployeeId:  employeeDto.EmployeeId,
		FirstName:   employeeDto.FirstName,
		LastName:    employeeDto.LastName,
		Email:       employeeDto.Email,
		PhoneNumber: employeeDto.PhoneNumber,
		IsAdmin:     isAdmin,
	}
	return employee, nil
}

func employeeDtoToEmployeeAuth(employeeDto EmployeeDto) domain.EmployeeAuth {
	employeeAuth := domain.EmployeeAuth{
		EmployeeId:   employeeDto.EmployeeId,
		Username:     employeeDto.Username,
		PasswordHash: employeeDto.Password,
	}
	return employeeAuth
}

func mapErrorsToEmployeeDto(errors domain.EmployeeErrors, employeeDto EmployeeDto) EmployeeDto {
	employeeDto.EmployeeIdErr = errors.EmployeeIdErr
	employeeDto.FirstNameErr = errors.FirstNameErr
	employeeDto.LastNameErr = errors.LastNameErr
	employeeDto.EmailErr = errors.EmailErr
	employeeDto.PhoneNumberErr = errors.PhoneNumberErr
	return employeeDto
}

func PostEmployee(employeeDto EmployeeDto, db *sql.DB) (EmployeeDto, error) {
	employee, err := employeeDtoToEmployee(employeeDto)
	if err != nil {
		return EmployeeDto{}, err
	}

	employeeAuth := employeeDtoToEmployeeAuth(employeeDto)

	errors := employee.Validate()
	if errors.IsSuccessful == false {
		employeeDto = mapErrorsToEmployeeDto(errors, employeeDto)
		return employeeDto, nil
	} else {
		log.Println(employee.PhoneNumber)
		employee, err = infrastructure.PostEmployee(employee, db)
		if err != nil {
			return EmployeeDto{}, err
		}
		employeeAuth.PasswordHash, err = auth.HashPassword(employeeAuth.PasswordHash)
		if err != nil {
			return EmployeeDto{}, err
		}
		employeeAuth, err = infrastructure.PostEmployeeAuth(employeeAuth, db)
		employeeDto.SuccessMsg = "Employee submitted successfully."
		return employeeDto, nil
	}
}

func returnEmployeeAndAuth(id, employeeId string, db *sql.DB) (domain.Employee, domain.EmployeeAuth, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return domain.Employee{}, domain.EmployeeAuth{}, err
	}
	outEmployee, err := infrastructure.GetEmployeeById(idInt, db)
	if err != nil {
		return outEmployee, domain.EmployeeAuth{}, err
	}
	outEmployeeAuth, err := infrastructure.GetEmployeeAuthByEmployeeId(employeeId, db)
	if err != nil {
		return outEmployee, outEmployeeAuth, err
	}

	return outEmployee, outEmployeeAuth, nil
}

func PutEmployee(employeeDto EmployeeDto, db *sql.DB) (EmployeeDto, error) {
	employee, err := employeeDtoToEmployee(employeeDto)
	if err != nil {
		return EmployeeDto{}, err
	}

	employeeAuth := employeeDtoToEmployeeAuth(employeeDto)

	errors := employee.Validate()
	if errors.IsSuccessful == false {
		employeeDto = mapErrorsToEmployeeDto(errors, employeeDto)
		return employeeDto, nil
	} else {
		employee, err = infrastructure.PutEmployee(employee, db)
		if err != nil {
			return EmployeeDto{}, err
		}
		employeeAuth.PasswordHash, err = auth.HashPassword(employeeAuth.PasswordHash)
		if err != nil {
			return EmployeeDto{}, err
		}
		_, err = infrastructure.PutEmployeeAuth(employeeAuth, db)
		employeeDto.SuccessMsg = "Employee submitted successfully."
		return employeeDto, nil
	}
}
