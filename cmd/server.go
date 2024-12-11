package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func newDb(sugar *zap.SugaredLogger) *sql.DB {
	env := os.Getenv("ENV")
	var dbPath string
	if env == "development" {
		dbPath = filepath.Join("..", "db", "dev.db")
	} else if env == "production" {
		dbPath = filepath.Join("..", "db", "prod.db")
	}
	if dbPath == "" {
		sugar.Error("error obtaining db path")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		sugar.Fatal(err)
	}

	return db
}

func new_server(ctx context.Context) http.Handler {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := logger.Sugar()
	defer logger.Sync()

	db := newDb(sugar)

	mux := http.NewServeMux()
	add_routes(mux, ctx, db, sugar)
	var handler http.Handler = mux
	return handler
}

func run(
	ctx context.Context,
) error {
	if err := godotenv.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "error loading .env file: %s\n", err)
	}
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	server := new_server(ctx)
	http_server := &http.Server{
		Addr:    "localhost:80",
		Handler: server,
	}
	go func() {
		fmt.Printf(" Listening and serving on %s\n", http_server.Addr)
		fmt.Println("Ctrl + C to exit")
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
