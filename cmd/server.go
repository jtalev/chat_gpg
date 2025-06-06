package main

import (
	"context"
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
	"github.com/jtalev/chat_gpg/domain/models"
	"github.com/jtalev/chat_gpg/handlers"
	"github.com/jtalev/chat_gpg/infrastructure/auth"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func new_server(
	ctx context.Context,
	h *handlers.Handler,
	a *infrastructure.Auth,
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
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rootPath, err := filepath.Abs("..")
	if err != nil {
		log.Fatalf("Error determinging root path: %v", err)
	}

	envPath := filepath.Join(rootPath, ".env")
	if err := godotenv.Load(envPath); err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := logger.Sugar()
	defer logger.Sync()

	db := domain.InitDb(rootPath, sugar)
	store := newCookieStore()

	h := handlers.Handler{
		DB:     db,
		Store:  store,
		Logger: sugar,
	}

	err = h.StartWorkers()
	if err != nil {
		return err
	}
	h.RegisterRepos()
	h.RegisterServices()

	a := infrastructure.Auth{
		Db:     db,
		Logger: sugar,
		Store:  store,
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	server := new_server(ctx, &h, &a, store, sugar)
	http_server := &http.Server{
		Addr:    ":8080",
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

func newCookieStore() *sessions.CookieStore {
	var (
		key = []byte(os.Getenv("SESSION_HMAC_KEY"))
		// encryptionKey = []byte(os.Getenv("SESSION_ENC_KEY"))
		store = sessions.NewCookieStore(key, nil)
	)

	env := os.Getenv("ENV")
	if env == "production" {
		store.Options = &sessions.Options{
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
	} else if env == "development" {
		store.Options = &sessions.Options{
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
	}

	return store
}
