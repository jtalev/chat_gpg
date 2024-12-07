package handlers

import "net/http"

type DashboardData struct {

}

func getDashboardData() []DashboardData {
	data := []DashboardData{
		{},
	}
	return data
}

func ServeDashboardView(w http.ResponseWriter, r *http.Request) {
	data := getDashboardData()
	component := "dashboard"
	title := "Dashboard - GPG"
	renderTemplate(w, component, title, data)
}
