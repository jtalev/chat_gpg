package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

type TimesheetData struct {
}

func getTimesheetData() []TimesheetData {
	data := []TimesheetData{
		{},
	}
	return data
}

func ServeTimesheetsView(store *sessions.CookieStore, sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data := getTimesheetData()
			component := "timesheets"
			title := "Timesheets - GPG"
			renderTemplate(w, r, store, component, title, data)
		},
	)
}
