package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/jtalev/chat_gpg/auth"
	"github.com/jtalev/chat_gpg/handlers"
	"go.uber.org/zap"
)

func ServeLoginView(w http.ResponseWriter, r *http.Request) {
	login_path := filepath.Join("..", "ui", "views", "login.html")
	tmpl := template.Must(template.ParseFiles(login_path))
	tmpl.Execute(w, nil)

}

func LoginHandler(db *sql.DB, sugar *zap.SugaredLogger) http.Handler {
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

			employee, err := handlers.GetEmployeeByEmployeeId(employeeAuth.EmployeeId)
			sugar.Info(employee)
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		},
	)
}
