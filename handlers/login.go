package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

func Serve_login(w http.ResponseWriter, r *http.Request) {
	login_path := filepath.Join("..", "ui", "views", "login.html")
	tmpl := template.Must(template.ParseFiles(login_path))
	tmpl.Execute(w, nil)

}

func authenticate_user(username, password string) (validation_result, error) {
	result := validation_result{Is_valid: true, Msg: "error"}
	if !result.Is_valid {
		err := fmt.Errorf("authenticate_user: invalid credentials")
		fmt.Println(err)
		return result, err
	}
	return result, nil
}

func Login_handler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	// get user from database

	// check if usernames and passwords match

	// if they do, redirect to dashboard
	result, err := authenticate_user(username, password)
	if err != nil {
		login_path := filepath.Join("..", "ui", "views", "login.html")
		tmpl := template.Must(template.ParseFiles(login_path))
		tmpl.Execute(w, map[string]interface{}{"ErrorMsg": result.Msg})
	} else {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}
