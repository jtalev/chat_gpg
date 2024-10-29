package main

import (
	"context"
	"net/http"

	"github.com/jtalev/chat_gpg/validation"
)

func add_routes(mux *http.ServeMux, ctx context.Context) {
	fileServer := http.FileServer(http.Dir("../ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.Handle("/", Serve_index())
	mux.Handle("/login", Serve_login())
	mux.Handle("/dashboard", Serve_dashboard())

	mux.Handle("/validate-login", validation.Handle_login_validation())
}