package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

type AdminData struct {
}

func getAdminData() []AdminData {
	data := []AdminData{
		{},
	}
	return data
}

func ServeAdminView(store *sessions.CookieStore, sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data := getAdminData()
			component := "admin"
			title := "Admin - GPG"
			renderTemplate(w, r, store, component, title, data)
		},
	)
}
