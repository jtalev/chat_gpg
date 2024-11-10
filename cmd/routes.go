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

	mux.HandleFunc("/authenticate-user", handlers.Login_handler)
}
