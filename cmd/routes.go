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
	mux.Handle("/logout", h.LogoutHandler())
	mux.Handle("/authenticate-user", h.LoginHandler())

	mux.HandleFunc("/error", h.ServeErrorView)

	mux.Handle("/dashboard", a.AuthMiddleware(h.ServeDashboardView()))
	mux.Handle("/jobs", a.AuthMiddleware(h.ServeJobsView()))
	mux.Handle("/timesheets", a.AuthMiddleware(h.ServeTimesheetsView()))
	mux.Handle("/leave", a.AuthMiddleware(h.ServeLeaveView()))
	mux.Handle("/admin", a.AuthMiddleware(h.ServeAdminView()))
	mux.Handle("/account", a.AuthMiddleware(h.ServeAccountView()))

	// leave requests
	mux.Handle("/leave/get", a.AuthMiddleware(h.GetLeaveRequests()))
	mux.Handle("/leave/get-by-employee", a.AuthMiddleware(h.GetLeaveRequestsByEmployee()))
	mux.Handle("/leave/post", a.AuthMiddleware(h.PostLeaveRequest()))
	mux.Handle("/leave/put", a.AuthMiddleware(h.PutLeaveRequest()))
	mux.Handle("/leave/delete", a.AuthMiddleware(h.DeleteLeaveRequest()))
}
