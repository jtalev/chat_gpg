package application

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/jtalev/chat_gpg/domain/models"
	models "github.com/jtalev/chat_gpg/domain/models"
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

type EmployeeDto struct {
	ID                  string
	EmployeeId          string
	FirstName           string
	LastName            string
	Username            string
	Password            string
	Email               string
	PhoneNumber         string
	Role                string
	IsAdmin             string
	EmployeeIdErr       string
	FirstNameErr        string
	LastNameErr         string
	UsernameErr         string
	PasswordErr         string
	EmailErr            string
	PhoneNumberErr      string
	RoleErr             string
	IsAdminErr          string
	IsSuccess           bool
	IsDifferentPassword bool
	SuccessMsg          string
}

func employeeDtoToEmployee(employeeDto EmployeeDto) (domain.Employee, error) {
	id, err := strconv.Atoi(employeeDto.ID)
	if err != nil {
		log.Println("strconv")
		return domain.Employee{}, err
	}
	var isAdmin bool
	if employeeDto.IsAdmin == "true" {
		isAdmin = true
	} else if employeeDto.IsAdmin == "false" {
		isAdmin = false
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

func employeeDtoToEmployeeAuth(employeeDto EmployeeDto) (domain.EmployeeAuth, error) {
	id, err := strconv.Atoi(employeeDto.ID)
	if err != nil {
		return domain.EmployeeAuth{}, err
	}
	employeeAuth := domain.EmployeeAuth{
		AuthId:       id,
		EmployeeId:   employeeDto.EmployeeId,
		Username:     employeeDto.Username,
		PasswordHash: employeeDto.Password,
	}
	return employeeAuth, nil
}

func employeeDtoToRole(employeeDto EmployeeDto, db *sql.DB) Role {
	r := Role{}
	role := models.Role{
		UUID:       uuid.New().String(),
		EmployeeId: employeeDto.EmployeeId,
		Role:       employeeDto.Role,
	}
	r.Role = role
	r.db = db
	return r
}

func mapErrorsToEmployeeDto(empErrors domain.EmployeeErrors, empAuthErrors domain.EmployeeAuthErrors, employeeDto EmployeeDto) EmployeeDto {
	employeeDto.EmployeeIdErr = empErrors.EmployeeIdErr
	employeeDto.FirstNameErr = empErrors.FirstNameErr
	employeeDto.LastNameErr = empErrors.LastNameErr
	employeeDto.EmailErr = empErrors.EmailErr
	employeeDto.PhoneNumberErr = empErrors.PhoneNumberErr
	employeeDto.UsernameErr = empAuthErrors.UsernameErr
	employeeDto.PasswordErr = empAuthErrors.PasswordErr
	employeeDto.IsAdminErr = empErrors.IsAdminErr
	return employeeDto
}

func PostEmployee(employeeDto EmployeeDto, db *sql.DB) (EmployeeDto, error) {
	employee, err := employeeDtoToEmployee(employeeDto)
	if err != nil {
		return EmployeeDto{}, err
	}

	employeeAuth, err := employeeDtoToEmployeeAuth(employeeDto)
	if err != nil {
		return employeeDto, err
	}

	empErrors := employee.Validate()
	empAuthErrors := employeeAuth.Validate()
	if empErrors.IsSuccessful == false || empAuthErrors.IsSuccessful == false {
		employeeDto = mapErrorsToEmployeeDto(empErrors, empAuthErrors, employeeDto)
		return employeeDto, nil
	} else {
		employee, err = infrastructure.PostEmployee(employee, db)
		if err != nil {
			return EmployeeDto{}, err
		}
		employeeAuth.PasswordHash, err = auth.HashPassword(employeeAuth.PasswordHash)
		if err != nil {
			return EmployeeDto{}, err
		}
		employeeAuth.AuthId = employee.ID
		employeeAuth, err = infrastructure.PostEmployeeAuth(employeeAuth, db)
		r := employeeDtoToRole(employeeDto, db)
		_, err = r.PostEmployeeRole()
		if err != nil {
			return employeeDto, err
		}
		employeeDto.SuccessMsg = "Employee submitted successfully."
		employeeDto.IsSuccess = true
		return employeeDto, nil
	}
}

func PutEmployee(employeeDto EmployeeDto, db *sql.DB) (EmployeeDto, error) {
	isSamePassword := false
	if employeeDto.Password == "" {
		authId, err := strconv.Atoi(employeeDto.ID)
		if err != nil {
			return employeeDto, err
		}
		employeeAuth, err := infrastructure.GetEmployeeAuthById(authId, db)
		if err != nil {
			return employeeDto, err
		}
		employeeDto.Password = employeeAuth.PasswordHash
		isSamePassword = true
	}

	employee, err := employeeDtoToEmployee(employeeDto)
	if err != nil {
		return EmployeeDto{}, err
	}

	employeeAuth, err := employeeDtoToEmployeeAuth(employeeDto)
	if err != nil {
		return employeeDto, err
	}
	empErrors := employee.Validate()
	empAuthErrors := employeeAuth.Validate()
	if empErrors.IsSuccessful == false || empAuthErrors.IsSuccessful == false {
		employeeDto = mapErrorsToEmployeeDto(empErrors, empAuthErrors, employeeDto)
		return employeeDto, nil
	} else {

		employee, err = infrastructure.PutEmployee(employee, db)
		if err != nil {
			return EmployeeDto{}, err
		}
		if isSamePassword == false {
			employeeDto.IsDifferentPassword = true
			employeeAuth.PasswordHash, err = auth.HashPassword(employeeAuth.PasswordHash)
			if err != nil {
				return EmployeeDto{}, err
			}
		}
		_, err = infrastructure.PutEmployeeAuth(employeeAuth, db)
		r := employeeDtoToRole(employeeDto, db)
		_, err = r.PutEmployeeRole()
		if err != nil {
			return employeeDto, err
		}
		employeeDto.SuccessMsg = "Employee submitted successfully."
		employeeDto.IsSuccess = true
		return employeeDto, nil
	}
}
