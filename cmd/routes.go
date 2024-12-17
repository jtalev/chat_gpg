package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/jtalev/chat_gpg/handlers"
	"go.uber.org/zap"
)

func add_routes(mux *http.ServeMux, ctx context.Context, db *sql.DB, sugar *zap.SugaredLogger) {
	fileServer := http.FileServer(http.Dir("../ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/login", handlers.ServeLoginView)
	mux.HandleFunc("/dashboard", handlers.ServeDashboardView)
	mux.HandleFunc("/jobs", handlers.ServeJobsView)
	mux.HandleFunc("/timesheets", handlers.ServeTimesheetsView)
	mux.Handle("/leave", handlers.ServeLeaveView(sugar))
	mux.HandleFunc("/admin", handlers.ServeAdminView)
	mux.HandleFunc("/account", handlers.ServeAccountView)

	// login requests
	mux.Handle("/authenticate-user", handlers.LoginHandler(db, sugar))

	// leave requests
	mux.Handle("/get-leave-requests", handlers.GetLeaveRequests(db, sugar))
	mux.Handle("/get-leave-request-by-id", handlers.GetLeaveRequestById(sugar))
	mux.Handle("/post-leave-request", handlers.PostLeaveRequest(sugar))
	mux.Handle("/put-leave-request", handlers.UpdateLeaveRequest(sugar))
	mux.Handle("/delete-leave-request", handlers.DeleteLeaveRequest(sugar))
}
