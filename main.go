package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Config from environment variables, with sensible defaults for local dev.
	// In production, set these via your hosting platform (Fly.io secrets, etc.)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://localhost:5432/url_shortener"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// Connect to the database. If this fails, crash immediately.
	// There's no point running a web server that can't reach its database.
	store, err := NewStore(databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer store.Close()

	// Chi is the HTTP router. It matches incoming requests to handler functions
	// based on the method (GET, POST, DELETE) and URL path.
	r := chi.NewRouter()

	// Middleware runs on every request before your handler.
	// Logger prints a log line for each request (method, path, status, duration).
	// Recoverer catches panics in handlers and returns a 500 instead of crashing the server.
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health check must be registered before the /{code} catch-all,
	// otherwise Chi would try to look up "health" as a short code.
	r.Get("/health", handleHealthCheck(store))

	// API routes (these serve JSON)
	r.Post("/api/shorten", handleShorten(store, baseURL))
	r.Get("/api/urls", handleListURLs(store))
	r.Get("/api/urls/{code}/stats", handleGetStats(store))
	r.Delete("/api/urls/{code}", handleDeleteURL(store))

	// Redirect route. {code} is a URL parameter, like a wildcard.
	// Any GET request that doesn't match /api/* or /health falls through to here.
	r.Get("/{code}", handleRedirect(store))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// http.Server with explicit timeouts instead of the zero-value defaults
	// from http.ListenAndServe (which waits forever on slow clients).
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start the server in a goroutine so we can listen for shutdown signals.
	// ListenAndServe blocks, so without the goroutine we'd never reach signal.Notify.
	go func() {
		fmt.Printf("Server starting on :%s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Block until we receive SIGINT (Ctrl+C) or SIGTERM (what Fly sends on deploy).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	// Give in-flight requests 10 seconds to finish (matches Fly's default kill_timeout).
	// After the deadline, Shutdown returns an error and the process exits.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	fmt.Println("Server stopped")
}
