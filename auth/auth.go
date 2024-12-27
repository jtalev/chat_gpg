package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jtalev/chat_gpg/models"
	"github.com/jtalev/chat_gpg/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return passwordHash, err
}

func VerifyHashedPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func AuthenticateUser(username, password string, db *sql.DB, sugar *zap.SugaredLogger) (models.EmployeeAuth, error) {
	employeeAuth, err := repository.GetEmployeeAuthByUsername(username, db)
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

type Auth struct {
	Logger *zap.SugaredLogger
	Store  *sessions.CookieStore
}

func (a *Auth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			session, err := a.Store.Get(r, "auth_session")
			if err != nil {
				a.Logger.Errorf("Error getting auth_session: %v", err)
				http.Error(w, "Error getting auth_session", http.StatusInternalServerError)
				return
			}
			if auth, ok := session.Values["is_authenticated"].(bool); !ok || !auth {
				a.Logger.Error("User is not authorized")
				http.Redirect(w, r, "/error", http.StatusFound)
				return
			}

			ctx := context.WithValue(r.Context(), "is_admin", session.Values["is_admin"])
			ctx = context.WithValue(ctx, "employee_id", session.Values["employee_id"])
			ctx = context.WithValue(ctx, "is_authenticated", session.Values["is_authenticated"])
			session.Save(r, w)

			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
