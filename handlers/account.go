package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

type AccountData struct {
}

func getAccountData() []AccountData {
	data := []AccountData{
		{},
	}
	return data
}

func ServeAccountView(store *sessions.CookieStore, sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data := getAccountData()
			component := "account"
			title := "Account - GPG"
			renderTemplate(w, r, store, component, title, data)
		},
	)
}
