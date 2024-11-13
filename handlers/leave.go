package handlers

import "net/http"

func Serve_leave(w http.ResponseWriter, r *http.Request) {
	component := "leave"
	title := "Leave - GPG"
	render_template(w, component, title)
}
