package handlers

import "net/http"

func Serve_admin(w http.ResponseWriter, r *http.Request) {
	component := "admin"
	title := "Admin - GPG"
	render_template(w, component, title)
}
