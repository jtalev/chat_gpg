package handlers

import (
	"fmt"
	"net/http"
)

func Serve_leave(w http.ResponseWriter, r *http.Request) {
	component := "leave"
	title := "Leave - GPG"
	render_template(w, component, title)
}

func Submit_leave_request(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "SUBMITTED")
}