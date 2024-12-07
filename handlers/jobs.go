package handlers

import "net/http"

type JobsData struct {

}

func getJobsData() []JobsData {
	data := []JobsData{
		{},
	}
	return data
}

func ServeJobsView(w http.ResponseWriter, r *http.Request) {
	data := getJobsData()
	component := "jobs"
	title := "Jobs - GPG"
	renderTemplate(w, component, title, data)
}
