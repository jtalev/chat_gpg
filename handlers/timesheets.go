package handlers

import "net/http"

func Serve_timesheets(w http.ResponseWriter, r *http.Request) {
	component := "timesheets"
	title := "Timesheets - GPG"
	render_template(w, component, title)
}
