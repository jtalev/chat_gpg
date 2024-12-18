package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/jtalev/chat_gpg/auth"
	"go.uber.org/zap"
)

func LoginHandler(db *sql.DB, store *sessions.CookieStore, sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			username := r.FormValue("username")
			password := r.FormValue("password")

			employeeAuth, err := auth.AuthenticateUser(username, password, db, sugar)
			if err != nil {
				login_path := filepath.Join("..", "ui", "views", "login.html")
				tmpl := template.Must(template.ParseFiles(login_path))
				tmpl.Execute(w, map[string]interface{}{"ErrorMsg": "Invalid username or password"})
				return
			}

			employee, err := GetEmployeeByEmployeeId(employeeAuth.EmployeeId, db)
			if err != nil {
				sugar.Errorf("Error getting authorised employee: %v", err)
				http.Error(w, "Error getting authorised employee", http.StatusInternalServerError)
				return
			}

			session, err := store.Get(r, "employee_session")
			if err != nil {
				sugar.Errorf("Error getting session", err)
				http.Error(w, "Error getting session", http.StatusInternalServerError)
				return
			}
			session.Values["is_authenticated"] = true
			session.Values["username"] = employee.EmployeeId
			if employee.IsAdmin {
				session.Values["is_admin"] = true
			} else {
				session.Values["is_admin"] = false
			}
			err = session.Save(r, w)
			if err != nil {
				sugar.Errorf("Error saving session: %v", err)
				http.Error(w, "Error saving session", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)
		},
	)
}

func LogoutHandler(store *sessions.CookieStore, sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// set is_authenticated cookie to false
			session, err := store.Get(r, "employee_session")
			if err != nil {
				sugar.Errorf("Error getting session: %v", err)
				http.Error(w, "Error getting session", http.StatusInternalServerError)
			}
			session.Values["is_authenticated"] = false

			// redirect user to login page
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		},
	)
}
