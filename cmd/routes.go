package main

import (
	"context"
	"net/http"

	"github.com/jtalev/chat_gpg/handlers"
)

func add_routes(mux *http.ServeMux, ctx context.Context) {
	fileServer := http.FileServer(http.Dir("../ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/login", handlers.ServeLoginView)
	mux.HandleFunc("/dashboard", handlers.ServeDashboardView)
	mux.HandleFunc("/jobs", handlers.ServeJobsView)
	mux.HandleFunc("/timesheets", handlers.ServeTimesheetsView)
	mux.HandleFunc("/leave", handlers.ServeLeaveView)
	mux.HandleFunc("/admin", handlers.ServeAdminView)
	mux.HandleFunc("/account", handlers.ServeAccountView)

	mux.HandleFunc("/authenticate-user", handlers.Login_handler)
	mux.HandleFunc("/submit-leave-request", handlers.Submit_leave_request)
}
