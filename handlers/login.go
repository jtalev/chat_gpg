package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func Serve_login(w http.ResponseWriter, r *http.Request) {
	login_path := filepath.Join("..", "ui", "views", "login.html")
	tmpl := template.Must(template.ParseFiles(login_path))
	tmpl.Execute(w, nil)
	
}

func validate_login(username, password string) validation_result {
	result := validation_result{Is_valid: true, Msg: ""}
	return result
}

func Login_handler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	// get user from database 
	
	// check if usernames and passwords match

	// if they do, redirect to dashboard
	result := validate_login(username, password)
	if result.Is_valid {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

	// if they don't, re-render page with error message dispalayed
}