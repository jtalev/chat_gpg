package handlers

import "net/http"

func Serve_jobs(w http.ResponseWriter, r *http.Request) {
	component := "jobs"
	title := "Jobs - GPG"
	render_template(w, component, title)
}
