package handlers

import (
	"net/http"
)

type TimesheetData struct {

}

func getTimesheetData() []TimesheetData {
	data := []TimesheetData{
		{},
	}
	return data
}

func ServeTimesheetsView(w http.ResponseWriter, r *http.Request) {
	data := getTimesheetData()
	component := "timesheets"
	title := "Timesheets - GPG"
	renderTemplate(w, component, title, data)
}
