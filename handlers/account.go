package handlers

import "net/http"

type AccountData struct {
}

func getAccountData() []AccountData {
	data := []AccountData{
		{},
	}
	return data
}

func ServeAccountView(w http.ResponseWriter, r *http.Request) {
	data := getAccountData()
	component := "account"
	title := "Account - GPG"
	renderTemplate(w, component, title, data)
}
