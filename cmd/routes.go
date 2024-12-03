package main

import (
	"context"
	"net/http"

	"github.com/jtalev/chat_gpg/handlers"
)

func add_routes(mux *http.ServeMux, ctx context.Context) {
	fileServer := http.FileServer(http.Dir("../ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/login", handlers.Serve_login)
	mux.HandleFunc("/dashboard", handlers.Serve_dashboard)
	mux.HandleFunc("/jobs", handlers.Serve_jobs)
	mux.HandleFunc("/timesheets", handlers.Serve_timesheets)
	mux.HandleFunc("/leave", handlers.Serve_leave)
	mux.HandleFunc("/admin", handlers.Serve_admin)
	mux.HandleFunc("/account", handlers.Serve_account)

	mux.HandleFunc("/authenticate-user", handlers.Login_handler)
	mux.HandleFunc("/submit-leave-request", handlers.Submit_leave_request)
}
