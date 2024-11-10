package handlers

import "net/http"

func Serve_dashboard(w http.ResponseWriter, r *http.Request) {
	component := "dashboard"
	title := "Dashboard"
	render_template(w, component, title)
}
