package auth

import (
	"database/sql"
	"errors"
	"log"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type EmployeeAuth struct {
	AuthId       int    `json:"auth_id"`
	EmployeeId   string `json:"employee_id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func GetEmployeeAuthByUsername(username string, db *sql.DB) (EmployeeAuth, error) {
	employeeAuth := EmployeeAuth{}
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

type ValidationResult struct {
	IsValid bool
	Msg     string
}

func HashPassword(password string) ([]byte, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return passwordHash, err
}

func VerifyHashedPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func AuthenticateUser(username, password string, db *sql.DB, sugar *zap.SugaredLogger) (EmployeeAuth, error) {
	employeeAuth, err := GetEmployeeAuthByUsername(username, db)
	if err == sql.ErrNoRows {
		sugar.Warnf("User %s not found", username)
		return employeeAuth, err
	} else if err != nil {
		sugar.Errorf("Database error: %v", err)
		return employeeAuth, err
	}
	if !VerifyHashedPassword(password, employeeAuth.PasswordHash) {
		sugar.Warnf("Invalid password for user %s", username)
		return employeeAuth, errors.New("invalid password")
	}
	return employeeAuth, nil
}
