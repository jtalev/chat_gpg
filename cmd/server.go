package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
	"path/filepath"
	"text/template"
)

// ----- output colours -----
var Yellow = "\033[33m"

func new_server(ctx context.Context) http.Handler {
	mux := http.NewServeMux()
	add_routes(mux, ctx)
	var handler http.Handler = mux
	return handler
}

func run(
	ctx context.Context,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	server := new_server(ctx)
	http_server := &http.Server{
		Addr:    "localhost:80",
		Handler: server,
	}
	go func() {
		fmt.Printf(" Listening and serving on %s\n", http_server.Addr)
		fmt.Println(string(Yellow), "Ctrl + C to exit")
		if err := http.ListenAndServe(http_server.Addr, http_server.Handler); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := http_server.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}

func Serve_index() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			path := filepath.Join("..", "ui", "pages", "index.html")
			tmpl := template.Must(template.ParseFiles(path))
			if err := tmpl.Execute(w, nil); err != nil {
				fmt.Printf("error executing template: %v", err)
				http.Error(w, "error executing tempate", http.StatusInternalServerError)
			}
		},
	)
}

func Serve_login() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			path := filepath.Join("..", "ui", "pages", "login.html")
			tmpl := template.Must(template.ParseFiles(path))
			if err := tmpl.Execute(w, nil); err != nil {
				fmt.Printf("error executing template: %v", err)
				http.Error(w, "error executing tempate", http.StatusInternalServerError)
			}
		},
	)
}

func Serve_dashboard() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			path := filepath.Join("..", "ui", "pages", "dashboard.html")
			nav_path := filepath.Join("..", "ui", "components", "nav.html")

			tmpl := template.Must(template.ParseFiles(path, nav_path))
			if err := tmpl.Execute(w, nil); err != nil {
				fmt.Printf("error executing template: %v", err)
				http.Error(w, "error executing tempate", http.StatusInternalServerError)
			}
		},
	)
}