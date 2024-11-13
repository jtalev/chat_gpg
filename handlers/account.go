package handlers

import "net/http"

func Serve_account(w http.ResponseWriter, r *http.Request) {
	component := "account"
	title := "Account - GPG"
	render_template(w, component, title)
}
