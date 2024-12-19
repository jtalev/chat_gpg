package main

import (
	"context"
	"net/http"

	"github.com/jtalev/chat_gpg/auth"
	"github.com/jtalev/chat_gpg/handlers"
)

func add_routes(mux *http.ServeMux, ctx context.Context, h *handlers.Handler, a *auth.Auth) {
	fileServer := http.FileServer(http.Dir("../ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	
	// login/logout requests
	mux.HandleFunc("/login", h.ServeLoginView)
	mux.Handle("/authenticate-user", h.LoginHandler())
	mux.Handle("/logout", h.LogoutHandler())

	mux.Handle("/dashboard", a.AuthMiddleware(h.ServeDashboardView()))
	mux.Handle("/jobs", a.AuthMiddleware(h.ServeJobsView()))
	mux.Handle("/timesheets", a.AuthMiddleware(h.ServeTimesheetsView()))
	mux.Handle("/leave", a.AuthMiddleware(h.ServeLeaveView()))
	mux.Handle("/admin", a.AuthMiddleware(h.ServeAdminView()))
	mux.Handle("/account", a.AuthMiddleware(h.ServeAccountView()))

	// leave requests
	mux.Handle("/post-leave-request", a.AuthMiddleware(handlers.PostLeaveRequest(h.DB, h.Sugar)))
}
