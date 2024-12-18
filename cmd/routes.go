package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jtalev/chat_gpg/handlers"
	"go.uber.org/zap"
)

func add_routes(mux *http.ServeMux, ctx context.Context, db *sql.DB, store *sessions.CookieStore, sugar *zap.SugaredLogger) {
	fileServer := http.FileServer(http.Dir("../ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/login", handlers.ServeLoginView)
	mux.Handle("/dashboard", handlers.ServeDashboardView(store, sugar))
	mux.Handle("/jobs", handlers.ServeJobsView(store, sugar))
	mux.Handle("/timesheets", handlers.ServeTimesheetsView(store, sugar))
	mux.Handle("/leave", handlers.ServeLeaveView(store, sugar))
	mux.Handle("/admin", handlers.ServeAdminView(store, sugar))
	mux.Handle("/account", handlers.ServeAccountView(store, sugar))

	// login/logout requests
	mux.Handle("/authenticate-user", handlers.LoginHandler(db, store, sugar))
	mux.Handle("/logout", handlers.LogoutHandler(store, sugar))

	// leave requests
	mux.Handle("/get-leave-requests", handlers.GetLeaveRequests(db, sugar))
	mux.Handle("/get-leave-request-by-id", handlers.GetLeaveRequestById(sugar))
	mux.Handle("/post-leave-request", handlers.PostLeaveRequest(sugar))
	mux.Handle("/put-leave-request", handlers.UpdateLeaveRequest(sugar))
	mux.Handle("/delete-leave-request", handlers.DeleteLeaveRequest(sugar))
}
