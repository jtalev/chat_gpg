package handlers

import "net/http"

type InitialDashboardData struct {
	Pages []string
}

func getDashboardData(w http.ResponseWriter, r *http.Request) (InitialDashboardData, error) {
	pages := []string{"TIMESHEETS", "LEAVE"}
	isAdmin, err := getIsAdmin(w, r)
	if err != nil {
		return InitialDashboardData{}, err
	}
	if isAdmin {
		pages = append(pages, "REPORTS", "ADMIN")
	}
	data := InitialDashboardData{
		Pages: pages,
	}
	return data, nil
}
