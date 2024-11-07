package handlers

import "net/http"

func Serve_dashboard(w http.ResponseWriter, r *http.Request) {
	render_template(w, "dashboard")
}