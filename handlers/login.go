package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/jtalev/chat_gpg/auth"
)

func (h *Handler) LoginHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			username := r.FormValue("username")
			password := r.FormValue("password")

			employeeAuth, err := auth.AuthenticateUser(username, password, h.DB, h.Sugar)
			if err != nil {
				login_path := filepath.Join("..", "ui", "views", "login.html")
				tmpl := template.Must(template.ParseFiles(login_path))
				tmpl.Execute(w, map[string]interface{}{"ErrorMsg": "Invalid username or password"})
				return
			}

			employee, err := GetEmployeeByEmployeeId(employeeAuth.EmployeeId, h.DB)
			if err != nil {
				h.Sugar.Errorf("Error getting authorised employee: %v", err)
				http.Error(w, "Error getting authorised employee", http.StatusInternalServerError)
				return
			}

			session, err := h.Store.Get(r, "auth_session")
			if err != nil {
				h.Sugar.Errorf("Error getting session", err)
				http.Error(w, "Error getting session", http.StatusInternalServerError)
				return
			}
			session.Values["is_authenticated"] = true
			session.Values["employee_id"] = employee.EmployeeId
			if employee.IsAdmin {
				session.Values["is_admin"] = true
			} else {
				session.Values["is_admin"] = false
			}
			err = session.Save(r, w)
			if err != nil {
				h.Sugar.Errorf("Error saving session: %v", err)
				http.Error(w, "Error saving session", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)
		},
	)
}

func (h *Handler) LogoutHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// set is_authenticated cookie to false
			session, err := h.Store.Get(r, "auth_session")
			if err != nil {
				h.Sugar.Errorf("Error getting session: %v", err)
				http.Error(w, "Error getting session", http.StatusInternalServerError)
			}
			session.Values["is_authenticated"] = false

			// redirect user to login page
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		},
	)
}
