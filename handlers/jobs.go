package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

type JobsData struct {
}

func getJobsData() []JobsData {
	data := []JobsData{
		{},
	}
	return data
}

func ServeJobsView(store *sessions.CookieStore, sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data := getJobsData()
			component := "jobs"
			title := "Jobs - GPG"
			renderTemplate(w, r, store, component, title, data)
		},
	)
}
