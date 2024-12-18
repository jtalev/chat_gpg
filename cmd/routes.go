package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jtalev/chat_gpg/auth"
	"github.com/jtalev/chat_gpg/handlers"
	"go.uber.org/zap"
)

func add_routes(mux *http.ServeMux, ctx context.Context, db *sql.DB, store *sessions.CookieStore, sugar *zap.SugaredLogger) {
	fileServer := http.FileServer(http.Dir("../ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/dashboard", auth.AuthMiddleware(handlers.ServeDashboardView(sugar), store, sugar))
	mux.Handle("/jobs", auth.AuthMiddleware(handlers.ServeJobsView(sugar), store, sugar))
	mux.Handle("/timesheets", auth.AuthMiddleware(handlers.ServeTimesheetsView(sugar), store, sugar))
	mux.Handle("/leave", auth.AuthMiddleware(handlers.ServeLeaveView(db, sugar), store, sugar))
	mux.Handle("/admin", auth.AuthMiddleware(handlers.ServeAdminView(sugar), store, sugar))
	mux.Handle("/account", auth.AuthMiddleware(handlers.ServeAccountView(sugar), store, sugar))

	// login/logout requests
	mux.HandleFunc("/login", handlers.ServeLoginView)
	mux.Handle("/authenticate-user", handlers.LoginHandler(db, store, sugar))
	mux.Handle("/logout", handlers.LogoutHandler(store, sugar))

	// leave requests
	// mux.Handle("get-leave-requests-by-employee-id", auth.AuthMiddleware(handlers.GetLeaveRequestsByEmployee(db, sugar), store, sugar))
	mux.Handle("/post-leave-request", auth.AuthMiddleware(handlers.PostLeaveRequest(db, sugar), store, sugar))
}
