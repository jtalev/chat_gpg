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

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/jtalev/chat_gpg/auth"
	"github.com/jtalev/chat_gpg/handlers"
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
		sugar.Error("Error obtaining db path")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		sugar.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		sugar.Error("DB connection not open:", err)
	} else {
		sugar.Info("DB connection is open and healthy")
	}

	return db
}

func newStore(sugar *zap.SugaredLogger) *sessions.CookieStore {
	var (
		key = []byte(os.Getenv("SESSION_HMAC_KEY"))
		// encryptionKey = []byte(os.Getenv("SESSION_ENC_KEY"))
		store = sessions.NewCookieStore(key, nil)
	)

	return store
}

func new_server(
	ctx context.Context,
	h *handlers.Handler,
	a *auth.Auth,
	store *sessions.CookieStore,
	sugar *zap.SugaredLogger,
) http.Handler {
	mux := http.NewServeMux()
	add_routes(mux, ctx, h, a)
	var handler http.Handler = mux
	return handler
}

func run(
	ctx context.Context,
) error {
	if err := godotenv.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %s\n", err)
	}
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := logger.Sugar()
	defer logger.Sync()
	db := newDb(sugar)
	store := newStore(sugar)
	h := handlers.Handler{
		DB: db,
		Store: store,
		Sugar: sugar,
	}
	a := auth.Auth{
		Sugar: sugar,
		Store: store,
	}
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	server := new_server(ctx, &h, &a, store, sugar)
	http_server := &http.Server{
		Addr:    "localhost:80",
		Handler: server,
	}
	go func() {
		fmt.Printf(" Listening and serving on %s\n", http_server.Addr)
		fmt.Println("Ctrl + C to exit")
		if err := http.ListenAndServe(http_server.Addr, http_server.Handler); err != nil && err != http.ErrServerClosed {
			sugar.Errorf("Error listening and serving: %s", err)
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
			fmt.Fprintf(os.Stderr, "Error shutting down http server: %s\n", err)
		}
		if err := db.Close(); err != nil {
			sugar.Error("Error closing DB connecion:", err)
		}
	}()
	wg.Wait()
	return nil
}
