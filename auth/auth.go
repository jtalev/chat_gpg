package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
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

func IsAuthenticated(r *http.Request, store *sessions.CookieStore) (bool, error) {
	session, err := store.Get(r, "employee_session")
	if err != nil {
		return false, err
	}
	cookie := session.Values["is_authenticated"]
	isAuthenticated := false
	if cookie == true {
		isAuthenticated = true
	}
	fmt.Println(isAuthenticated)
	return isAuthenticated, nil
}

func RedirectUnauthorisedUser(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore, sugar *zap.SugaredLogger) {
	isAuthenticated, err := IsAuthenticated(r, store)
	if err != nil {
		sugar.Errorf("Error getting authentication status: %v", err)
		http.Error(w, "Error getting authentication status", http.StatusInternalServerError)
		return
	}
	if isAuthenticated == false {
		sugar.Info("User is not authenticated")
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}
}