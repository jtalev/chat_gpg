package handlers

import "net/http"

type AdminData struct {
}

func getAdminData() []AdminData {
	data := []AdminData{
		{},
	}
	return data
}

func ServeAdminView(w http.ResponseWriter, r *http.Request) {
	data := getAdminData()
	component := "admin"
	title := "Admin - GPG"
	renderTemplate(w, component, title, data)
}
