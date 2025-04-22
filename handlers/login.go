package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	auth "github.com/jtalev/chat_gpg/infrastructure/auth"
	repository "github.com/jtalev/chat_gpg/infrastructure/repository"
)

func (h *Handler) LoginHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			var username, password string
			if r.Header.Get("Content-Type") == "application/json" {
				var loginData struct {
					Username string `json:"username"`
					Password string `json:"password"`
				}

				if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
					http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
					return
				}

				username = loginData.Username
				password = loginData.Password
			} else {
				username = r.FormValue("username")
				password = r.FormValue("password")
			}

			employeeAuth, err := auth.AuthenticateUser(username, password, h.DB, h.Logger)
			if err != nil {
				login_path := filepath.Join("..", "ui", "views", "login.html")
				tmpl := template.Must(template.ParseFiles(login_path))
				tmpl.Execute(w, map[string]interface{}{"ErrorMsg": "Invalid username or password"})
				return
			}

			employee, err := repository.GetEmployeeByEmployeeId(employeeAuth.EmployeeId, h.DB)
			if err != nil {
				h.Logger.Errorf("Error getting authorised employee: %v", err)
				http.Error(w, "Error getting authorised employee", http.StatusInternalServerError)
				return
			}

			session, err := h.Store.Get(r, "auth_session")
			if err != nil {
				h.Logger.Errorf("Error getting session", err)
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
				h.Logger.Errorf("Error saving session: %v", err)
				http.Error(w, "Error saving session", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		},
	)
}

func (h *Handler) LogoutHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			h.Logger.Info("handling logout")

			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")

			session, err := h.Store.Get(r, "auth_session")
			if err != nil {
				h.Logger.Errorf("Error getting session: %v", err)
				http.Error(w, "Error getting session", http.StatusInternalServerError)
			}
			session.Values["is_authenticated"] = false
			session.Values["is_admin"] = false
			session.Values["employee_id"] = -1
			session.Options.MaxAge = -1
			err = session.Save(r, w)
			if err != nil {
				h.Logger.Errorf("Error saving session: %v", err)
				http.Error(w, "Error saving session", http.StatusInternalServerError)
				return
			}

			// redirect user to login page
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		},
	)
}
