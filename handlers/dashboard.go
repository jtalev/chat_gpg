package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

type DashboardData struct {
	IsAdmin bool
}

func getDashboardData() []DashboardData {
	data := []DashboardData{
		{},
	}
	return data
}

func ServeDashboardView(store *sessions.CookieStore, sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data := getDashboardData()
			component := "dashboard"
			title := "Dashboard - GPG"
			renderTemplate(w, r, store, component, title, data)
		},
	)
}
